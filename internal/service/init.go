package service

import (
	"context"

	api "github.com/cicbyte/itfeeds/api/v1/init"
)

type IInit interface {
	Status(ctx context.Context) (*api.InitStatusRes, error)
	TestConnection(ctx context.Context, req *api.InitTestReq) (*api.InitTestRes, error)
	Setup(ctx context.Context, req *api.InitSetupReq) (*api.InitSetupRes, error)
}

var localInit IInit

func Init() IInit {
	if localInit == nil {
		panic("implement not found for interface IInit, forgot register?")
	}
	return localInit
}

func RegisterInit(i IInit) {
	localInit = i
}
