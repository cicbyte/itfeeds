package controller

import (
	"context"
	"fmt"

	api "github.com/cicbyte/itfeeds/api/v1/rss_entries"
	consts "github.com/cicbyte/itfeeds/internal/consts"
	service "github.com/cicbyte/itfeeds/internal/service"
)

var RssEntries = rssEntriesController{}

type rssEntriesController struct {
	BaseController
}

// Add 新增RSS条目
func (c *rssEntriesController) Add(ctx context.Context, req *api.RssEntriesAddReq) (res *api.RssEntriesAddRes, err error) {
	res = new(api.RssEntriesAddRes)
	err = service.RssEntries().Add(ctx, req)
	return
}

// List 获取RSS条目列表
func (c *rssEntriesController) List(ctx context.Context, req *api.RssEntriesListReq) (res *api.RssEntriesListRes, err error) {
	res = new(api.RssEntriesListRes)
	if req.PageSize == 0 {
		req.PageSize = consts.PageSize
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	total, list, err := service.RssEntries().List(ctx, req)
	res.Total = total
	res.CurrentPage = req.PageNum
	res.List = list
	return
}

// Get 获取RSS条目详情
func (c *rssEntriesController) Get(ctx context.Context, req *api.RssEntriesDetailReq) (res *api.RssEntriesDetailRes, err error) {
	res = new(api.RssEntriesDetailRes)
	res.RssEntriesInfo, err = service.RssEntries().GetById(ctx, req.Id)
	return
}

// Edit 修改RSS条目
func (c *rssEntriesController) Edit(ctx context.Context, req *api.RssEntriesEditReq) (res *api.RssEntriesEditRes, err error) {
	res = new(api.RssEntriesEditRes)
	err = service.RssEntries().Edit(ctx, req)
	return
}

// Delete 删除RSS条目
func (c *rssEntriesController) Delete(ctx context.Context, req *api.RssEntriesDelReq) (res *api.RssEntriesDelRes, err error) {
	res = new(api.RssEntriesDelRes)
	err = service.RssEntries().Delete(ctx, req.Id)
	return
}

// BatchDelete 批量删除RSS条目
func (c *rssEntriesController) BatchDelete(ctx context.Context, req *api.RssEntriesBatchDelReq) (res *api.RssEntriesBatchDelRes, err error) {
	res = new(api.RssEntriesBatchDelRes)
	err = service.RssEntries().BatchDelete(ctx, req.Ids)
	return
}

// PublicList 公开RSS条目列表
func (c *rssEntriesController) PublicList(ctx context.Context, req *api.PublicRssEntriesListReq) (res *api.PublicRssEntriesListRes, err error) {
	res = new(api.PublicRssEntriesListRes)
	if req.PageSize == 0 {
		req.PageSize = consts.PageSize
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	listReq := &api.RssEntriesListReq{
		PageReq: req.PageReq,
		Title:   req.Title,
		Author:  req.Author,
	}
	total, list, err := service.RssEntries().List(ctx, listReq)
	res.Total = total
	res.CurrentPage = req.PageNum
	res.List = list
	return
}

// Sync 手动拉取RSS
func (c *rssEntriesController) Sync(ctx context.Context, req *api.RssEntriesSyncReq) (res *api.RssEntriesSyncRes, err error) {
	res = new(api.RssEntriesSyncRes)
	added, err := service.RssEntries().Sync(ctx)
	if err != nil {
		return
	}
	res.Added = added
	res.Message = fmt.Sprintf("新增 %d 条", added)
	return
}
