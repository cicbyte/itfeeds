package rss_entries

import (
	commonApi "github.com/cicbyte/itfeeds/api/v1/common"
	model "github.com/cicbyte/itfeeds/internal/model"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RssEntriesAddReq RSS条目-新增请求
type RssEntriesAddReq struct {
	g.Meta    `path:"/rss_entries/add" method:"post" tags:"RSS条目" summary:"RSS条目-新增"`
	Guid      string      `json:"guid" v:"required#GUID不能为空"`
	Url       string      `json:"url" v:"required#URL不能为空"`
	Title     string      `json:"title" v:"required#标题不能为空"`
	Content   string      `json:"content"`
	Published *gtime.Time `json:"published"`
	Author    string      `json:"author"`
}

// RssEntriesAddRes RSS条目-新增响应
type RssEntriesAddRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

// RssEntriesDelReq RSS条目-删除请求
type RssEntriesDelReq struct {
	g.Meta `path:"/rss_entries/del" method:"delete" tags:"RSS条目" summary:"RSS条目-删除"`
	Id     int `json:"id" v:"required#id不能为空"`
}

// RssEntriesDelRes RSS条目-删除响应
type RssEntriesDelRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

// RssEntriesBatchDelReq RSS条目-批量删除请求
type RssEntriesBatchDelReq struct {
	g.Meta `path:"/rss_entries/batchdel" method:"delete" tags:"RSS条目" summary:"RSS条目-批量删除"`
	Ids    []int `json:"ids" v:"required#ids不能为空"`
}

// RssEntriesBatchDelRes RSS条目-批量删除响应
type RssEntriesBatchDelRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

// RssEntriesEditReq RSS条目-修改请求
type RssEntriesEditReq struct {
	g.Meta    `path:"/rss_entries/edit" method:"put" tags:"RSS条目" summary:"RSS条目-修改"`
	Id        int         `json:"id" v:"required#ID不能为空"`
	Guid      string      `json:"guid" v:"required#GUID不能为空"`
	Url       string      `json:"url" v:"required#URL不能为空"`
	Title     string      `json:"title" v:"required#标题不能为空"`
	Content   string      `json:"content"`
	Published *gtime.Time `json:"published"`
	Author    string      `json:"author"`
}

// RssEntriesEditRes RSS条目-修改响应
type RssEntriesEditRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

// RssEntriesListReq RSS条目-列表请求
type RssEntriesListReq struct {
	g.Meta    `path:"/rss_entries/list" method:"get" tags:"RSS条目" summary:"RSS条目-列表"`
	commonApi.PageReq
	Title     string `json:"title" dc:"标题筛选"`
	Author    string `json:"author" dc:"作者筛选"`
	StartDate string `json:"startDate" dc:"开始日期(YYYY-MM-DD)"`
	EndDate   string `json:"endDate" dc:"结束日期(YYYY-MM-DD)"`
}

// RssEntriesListRes RSS条目-列表响应
type RssEntriesListRes struct {
	g.Meta        `mime:"application/json" example:"string"`
	commonApi.ListRes
	List []*model.RssEntriesInfo `json:"list"`
}

// RssEntriesDetailReq RSS条目-详情请求
type RssEntriesDetailReq struct {
	g.Meta `path:"/rss_entries/detail" method:"get" tags:"RSS条目" summary:"RSS条目-详情"`
	Id     int `json:"id" v:"required#id不能为空"`
}

// RssEntriesDetailRes RSS条目-详情响应
type RssEntriesDetailRes struct {
	g.Meta `mime:"application/json" example:"string"`
	*model.RssEntriesInfo
}

// RssEntriesSyncReq 手动拉取RSS请求
type RssEntriesSyncReq struct {
	g.Meta `path:"/rss_entries/sync" method:"post" tags:"RSS条目" summary:"手动拉取RSS"`
}

// RssEntriesSyncRes 手动拉取RSS响应
type RssEntriesSyncRes struct {
	g.Meta  `mime:"application/json"`
	Added   int    `json:"added"`
	Message string `json:"message,omitempty"`
}

// PublicRssEntriesListReq 获取公开RSS条目列表请求
type PublicRssEntriesListReq struct {
	g.Meta    `path:"/rss_entries" method:"get" tags:"公开接口" summary:"获取RSS条目列表"`
	commonApi.PageReq
	Title     string `json:"title" dc:"标题筛选"`
	Author    string `json:"author" dc:"作者筛选"`
	StartDate string `json:"startDate" dc:"开始日期(YYYY-MM-DD)"`
	EndDate   string `json:"endDate" dc:"结束日期(YYYY-MM-DD)"`
}

// PublicRssEntriesListRes 获取公开RSS条目列表响应
type PublicRssEntriesListRes struct {
	g.Meta        `mime:"application/json" example:"string"`
	commonApi.ListRes
	List []*model.RssEntriesInfo `json:"list"`
}
