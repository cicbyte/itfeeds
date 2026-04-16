package main

import (
	"github.com/cicbyte/itfeeds/internal/cmd"
	_ "github.com/cicbyte/itfeeds/internal/packed"
	//重要 需要导入数据库驱动
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

func main() {
	// 确保 config.yaml 存在，不存在时从 example 复制（首次部署场景）
	configPath := "manifest/config/config.yaml"
	examplePath := "manifest/config/config.example.yaml"
	if !gfile.Exists(configPath) && gfile.Exists(examplePath) {
		gfile.CopyFile(examplePath, configPath)
	}
	cmd.Main.Run(gctx.GetInitCtx())
}