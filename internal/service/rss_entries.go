package service

import (
	"context"

	api "github.com/cicbyte/itfeeds/api/v1/rss_entries"
	model "github.com/cicbyte/itfeeds/internal/model"
)

type IRssEntries interface {
	List(ctx context.Context, req *api.RssEntriesListReq) (total interface{}, res []*model.RssEntriesInfo, err error)
	Add(ctx context.Context, req *api.RssEntriesAddReq) (err error)
	Edit(ctx context.Context, req *api.RssEntriesEditReq) (err error)
	Delete(ctx context.Context, id int) (err error)
	BatchDelete(ctx context.Context, ids []int) (err error)
	GetById(ctx context.Context, id int) (res *model.RssEntriesInfo, err error)
	Sync(ctx context.Context) (added int, err error)
}

var localRssEntries IRssEntries

func RssEntries() IRssEntries {
	if localRssEntries == nil {
		panic("implement not found for interface IRssEntries, forgot register?")
	}
	return localRssEntries
}

func RegisterRssEntries(i IRssEntries) {
	localRssEntries = i
}
