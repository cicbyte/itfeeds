package controller

import (
	"context"

	api "github.com/cicbyte/itfeeds/api/v1/init"
	"github.com/cicbyte/itfeeds/internal/service"
)

var Init = initController{}

type initController struct {
	BaseController
}

func (c *initController) Status(ctx context.Context, req *api.InitStatusReq) (res *api.InitStatusRes, err error) {
	return service.Init().Status(ctx)
}

func (c *initController) TestConnection(ctx context.Context, req *api.InitTestReq) (res *api.InitTestRes, err error) {
	return service.Init().TestConnection(ctx, req)
}

func (c *initController) Setup(ctx context.Context, req *api.InitSetupReq) (res *api.InitSetupRes, err error) {
	return service.Init().Setup(ctx, req)
}
