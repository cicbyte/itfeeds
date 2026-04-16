package router

import (
	"context"

	controller "github.com/cicbyte/itfeeds/internal/controller"

	"github.com/cicbyte/itfeeds/internal/service"
	"github.com/cicbyte/itfeeds/library/libRouter"
	"github.com/gogf/gf/v2/net/ghttp"
)

type Router struct{}

func (router *Router) BindController(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.Middleware(service.Middleware().MiddlewareCORS)

		group.Bind(
			controller.Health,
			controller.RssEntries,
				controller.Init,
		)

		//自动绑定定义的控制器
		if err := libRouter.RouterAutoBind(ctx, router, group); err != nil {
			panic(err)
		}
	})
}