package init

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"time"

	_ "github.com/go-sql-driver/mysql"

	api "github.com/cicbyte/itfeeds/api/v1/init"
	"github.com/cicbyte/itfeeds/internal/logic/initdb"
	"github.com/cicbyte/itfeeds/internal/logic/rss_sync"
	"github.com/cicbyte/itfeeds/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
)

const initFlagPath = "manifest/config/.initialized"
const initializingFlagPath = "manifest/config/.initializing"

var validDBName = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func init() {
	service.RegisterInit(New())
}

func New() *sInit {
	return &sInit{}
}

type sInit struct{}

func (s *sInit) Status(ctx context.Context) (*api.InitStatusRes, error) {
	return &api.InitStatusRes{
		Initialized: gfile.Exists(initFlagPath),
	}, nil
}

func (s *sInit) TestConnection(ctx context.Context, req *api.InitTestReq) (*api.InitTestRes, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=true&loc=Local",
		req.User, req.Password, req.Host, req.Port)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return &api.InitTestRes{Success: false, Error: "创建连接失败: " + err.Error()}, nil
	}
	defer db.Close()

	var version string
	if err := db.QueryRow("SELECT VERSION()").Scan(&version); err != nil {
		return &api.InitTestRes{Success: false, Error: "连接失败: " + err.Error()}, nil
	}

	return &api.InitTestRes{Success: true, Version: gconv.String(version)}, nil
}

func (s *sInit) Setup(ctx context.Context, req *api.InitSetupReq) (*api.InitSetupRes, error) {
	// 检查是否已初始化
	if gfile.Exists(initFlagPath) {
		return &api.InitSetupRes{Success: false, Error: "系统已完成初始化"}, nil
	}

	// 检查是否正在初始化
	if gfile.Exists(initializingFlagPath) {
		return &api.InitSetupRes{Success: false, Error: "正在初始化中，请稍后重试"}, nil
	}

	// 校验数据库名
	if !validDBName.MatchString(req.Database) {
		return &api.InitSetupRes{Success: false, Error: "数据库名只允许字母、数字和下划线"}, nil
	}

	// 创建初始化锁
	gfile.PutContents(initializingFlagPath, time.Now().Format("2006-01-02 15:04:05"))
	cleanup := func() {
		if gfile.Exists(initializingFlagPath) {
			gfile.Remove(initializingFlagPath)
		}
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=true&loc=Local",
		req.User, req.Password, req.Host, req.Port)

	// 1. 连接 MySQL，创建数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		cleanup()
		return &api.InitSetupRes{Success: false, Error: "创建连接失败: " + err.Error()}, nil
	}
	defer db.Close()

	createDbSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci", req.Database)
	if _, err := db.Exec(createDbSQL); err != nil {
		cleanup()
		return &api.InitSetupRes{Success: false, Error: "创建数据库失败: " + err.Error()}, nil
	}

	// 2. 连接目标数据库，建表
	targetDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		req.User, req.Password, req.Host, req.Port, req.Database)
	targetDb, err := sql.Open("mysql", targetDSN)
	if err != nil {
		cleanup()
		return &api.InitSetupRes{Success: false, Error: "连接数据库失败: " + err.Error()}, nil
	}
	defer targetDb.Close()

	if _, err := targetDb.Exec(initdb.CreateTableSQL); err != nil {
		cleanup()
		return &api.InitSetupRes{Success: false, Error: "创建表失败: " + err.Error()}, nil
	}

	// 3. 写入 config.yaml
	link := fmt.Sprintf("mysql:%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Asia%%2FShanghai", req.User, req.Password, req.Host, req.Port, req.Database)
	if err := s.writeConfig(link); err != nil {
		cleanup()
		return &api.InitSetupRes{Success: false, Error: "写入配置失败: " + err.Error()}, nil
	}

	// 4. 创建标记文件
	gfile.PutContents(initFlagPath, gtime.Now().Format("Y-m-d H:i:s"))
	cleanup()

	// 5. 延迟启动 RSS 同步
	go func() {
		time.Sleep(2 * time.Second)
		if err := rss_sync.StartRSSSync(gctx.New()); err != nil {
			g.Log().Errorf(ctx, "启动RSS同步失败: %v", err)
		}
	}()

	return &api.InitSetupRes{Success: true, Message: "初始化成功"}, nil
}

func (s *sInit) writeConfig(link string) error {
	configPath := "manifest/config/config.yaml"
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	pattern := `(database:\s*\n\s*default:\s*\n\s*link:\s*")[^"]*(")`
	re := regexp.MustCompile(pattern)
	result := re.ReplaceAllString(string(content), "${1}"+link+"${2}")
	return os.WriteFile(configPath, []byte(result), 0644)
}
