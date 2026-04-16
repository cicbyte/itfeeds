package model

import "github.com/gogf/gf/v2/os/gtime"

// RssEntriesInfo RSS条目信息模型
type RssEntriesInfo struct {
	Id        int         `orm:"id" json:"id"`               // 主键ID，自增
	Guid      string      `orm:"guid" json:"guid"`           // 全局唯一标识符
	Url       string      `orm:"url" json:"url"`             // 文章URL
	Title     string      `orm:"title" json:"title"`         // 文章标题
	Content   string      `orm:"content" json:"content"`     // 文章内容
	Published *gtime.Time `orm:"published" json:"published"` // 发布时间
	Author    string      `orm:"author" json:"author"`       // 作者
}
