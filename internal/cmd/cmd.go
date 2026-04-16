package cmd

import (
	"context"
	"strings"

	"github.com/cicbyte/itfeeds/internal/controller"
	_ "github.com/cicbyte/itfeeds/internal/logic"
	"github.com/cicbyte/itfeeds/internal/logic/initdb"
	"github.com/cicbyte/itfeeds/internal/logic/rss_sync"
	"github.com/cicbyte/itfeeds/internal/mcp"
	"github.com/cicbyte/itfeeds/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()

			initialized := gfile.Exists("manifest/config/.initialized")

			if initialized {
				if err := initdb.EnsureTables(ctx); err != nil {
					g.Log().Errorf(ctx, "初始化数据库表失败: %v", err)
				}
				if err := rss_sync.StartRSSSync(ctx); err != nil {
					g.Log().Errorf(ctx, "启动RSS同步失败: %v", err)
				}
			} else {
				g.Log().Info(ctx, "系统未初始化，请通过 Web 页面完成初始化配置")
			}

			mcpHandler := mcp.NewStreamableHTTPServer()

			s.BindHandler("/mcp/*", func(r *ghttp.Request) {
				mcpHandler.ServeHTTP(r.Response.Writer.ResponseWriter, r.Request)
			})

			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)

				group.Group("/api/v1", func(apiGroup *ghttp.RouterGroup) {
					apiGroup.Middleware(service.Middleware().MiddlewareCORS)

					apiGroup.Bind(
						controller.Init,
						controller.Health,
						controller.RssEntries,
					)
				})

				// SPA 回退
				group.Hook("/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
					path := r.URL.Path

					if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/mcp") {
						return
					}

					if strings.Contains(path, ".") && !strings.HasSuffix(path, "/") {
						return
					}

					if path != "/" {
						r.Response.ServeFile("resource/public/index.html")
						r.ExitAll()
					}
				})
			})
			s.Run()
			return nil
		},
	}
)
