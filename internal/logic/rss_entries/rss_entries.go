package rss_entries

import (
	"context"
	"fmt"

	api "github.com/cicbyte/itfeeds/api/v1/rss_entries"
	dao "github.com/cicbyte/itfeeds/internal/dao"
	"github.com/cicbyte/itfeeds/internal/logic/rss_sync"
	model "github.com/cicbyte/itfeeds/internal/model"
	do "github.com/cicbyte/itfeeds/internal/model/do"
	service "github.com/cicbyte/itfeeds/internal/service"
	liberr "github.com/cicbyte/itfeeds/library/liberr"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

func init() {
	service.RegisterRssEntries(New())
}

func New() *sRssEntries {
	return &sRssEntries{}
}

type sRssEntries struct{}

// List 获取RSS条目列表
func (s sRssEntries) List(ctx context.Context, req *api.RssEntriesListReq) (total interface{}, list []*model.RssEntriesInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		m := dao.RssEntries.Ctx(ctx)
		columns := dao.RssEntries.Columns()

		// 标题模糊查询
		if req.Title != "" {
			m = m.Where(fmt.Sprintf("%s like ?", columns.Title), "%"+req.Title+"%")
		}
		// 作者模糊查询
		if req.Author != "" {
			m = m.Where(fmt.Sprintf("%s like ?", columns.Author), "%"+req.Author+"%")
		}
		// 时间范围查询
		if req.StartDate != "" {
			startTime := gtime.NewFromStrFormat(req.StartDate, "Y-m-d")
			if startTime != nil {
				m = m.Where(fmt.Sprintf("%s >= ?", columns.Published), startTime)
			}
		}
		if req.EndDate != "" {
			endTime := gtime.NewFromStrFormat(req.EndDate+" 23:59:59", "Y-m-d H:i:s")
			if endTime != nil {
				m = m.Where(fmt.Sprintf("%s <= ?", columns.Published), endTime)
			}
		}

		orderBy := req.OrderBy
		if orderBy == "" {
			orderBy = fmt.Sprintf("%s desc", columns.Published)
		}

		m = m.Safe()

		// 先 Count 再 Scan，Count 不需要排序，效率更高
		total, err = m.Count()
		liberr.ErrIsNil(ctx, err, "获取RSS条目总数失败")

		err = m.Page(req.PageNum, req.PageSize).Order(orderBy).Scan(&list)
		liberr.ErrIsNil(ctx, err, "获取RSS条目列表失败")
	})
	return
}

// Add 新增RSS条目
func (s sRssEntries) Add(ctx context.Context, req *api.RssEntriesAddReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = dao.RssEntries.Ctx(ctx).Insert(do.RssEntries{
			Guid:      req.Guid,
			Url:       req.Url,
			Title:     req.Title,
			Content:   req.Content,
			Published: req.Published,
			Author:    req.Author,
		})
		liberr.ErrIsNil(ctx, err, "新增RSS条目失败")
	})
	return
}

// Edit 修改RSS条目
func (s sRssEntries) Edit(ctx context.Context, req *api.RssEntriesEditReq) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = s.GetById(ctx, req.Id)
		liberr.ErrIsNil(ctx, err, "获取RSS条目失败")

		_, err = dao.RssEntries.Ctx(ctx).WherePri(req.Id).Data(do.RssEntries{
			Guid:      req.Guid,
			Url:       req.Url,
			Title:     req.Title,
			Content:   req.Content,
			Published: req.Published,
			Author:    req.Author,
		}).Update()
		liberr.ErrIsNil(ctx, err, "修改RSS条目失败")
	})
	return
}

// Delete 删除RSS条目
func (s sRssEntries) Delete(ctx context.Context, id int) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = dao.RssEntries.Ctx(ctx).WherePri(id).Delete()
		liberr.ErrIsNil(ctx, err, "删除RSS条目失败")
	})
	return
}

// BatchDelete 批量删除RSS条目
func (s sRssEntries) BatchDelete(ctx context.Context, ids []int) (err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		_, err = dao.RssEntries.Ctx(ctx).Where(dao.RssEntries.Columns().Id+" in(?)", ids).Delete()
		liberr.ErrIsNil(ctx, err, "批量删除RSS条目失败")
	})
	return
}

// GetById 根据ID获取RSS条目
func (s sRssEntries) GetById(ctx context.Context, id int) (res *model.RssEntriesInfo, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		err = dao.RssEntries.Ctx(ctx).Where(fmt.Sprintf("%s=?", dao.RssEntries.Columns().Id), id).Scan(&res)
		liberr.ErrIsNil(ctx, err, "获取RSS条目失败")
	})
	return
}

// Sync 手动拉取RSS
func (s sRssEntries) Sync(ctx context.Context) (added int, err error) {
	added = rss_sync.SyncAllFeeds(ctx)
	return
}
