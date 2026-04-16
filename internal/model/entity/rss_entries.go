// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RssEntries is the golang structure for table rss_entries.
type RssEntries struct {
	Id        int         `json:"id"        orm:"id"        description:""` //
	Guid      string      `json:"guid"      orm:"guid"      description:""` //
	Url       string      `json:"url"       orm:"url"       description:""` //
	Title     string      `json:"title"     orm:"title"     description:""` //
	Content   string      `json:"content"   orm:"content"   description:""` //
	Published *gtime.Time `json:"published" orm:"published" description:""` //
	Author    string      `json:"author"    orm:"author"    description:""` //
}
