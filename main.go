package main

import (
	"github.com/cicbyte/itfeeds/internal/cmd"
	_ "github.com/cicbyte/itfeeds/internal/packed"
	//重要 需要导入数据库驱动
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

// defaultConfig 首次部署时生成的默认配置（数据库连接为占位符，需通过 /init 页面配置）
const defaultConfig = `server:
  address: ":8000"
  serverRoot: "resource/public"
  logPath: "resource/log/server"
  logStdout: true
  errorStack: true
  errorLogEnabled: true
  errorLogPattern: "error-{Ymd}.log"
  accessLogEnabled: true
  accessLogPattern: "access-{Ymd}.log"

logger:
  path: "resource/log/run"
  file: "{Y-m-d}.log"
  level: "all"
  stdout: true

database:
  default:
    link: "mysql:root:@tcp(127.0.0.1:3306)/itfeeds?charset=utf8mb4&parseTime=true&loc=Asia%%2FShanghai"

rss:
  enabled: true
  crons:
    - "0 0 8-21 * * *"
    - "0 30 8-20 * * *"
  feeds:
    - "https://www.ithome.com/rss"
  barkPush: ""
`

func main() {
	configPath := "manifest/config/config.yaml"
	if !gfile.Exists(configPath) {
		gfile.Mkdir("manifest/config")
		gfile.PutContents(configPath, defaultConfig)
	}
	cmd.Main.Run(gctx.GetInitCtx())
}
