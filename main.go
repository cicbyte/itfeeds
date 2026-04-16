package main

import (
	_ "github.com/cicbyte/itfeeds/internal/packed"
	//重要 需要导入数据库驱动
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/os/gctx"

	"github.com/cicbyte/itfeeds/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}