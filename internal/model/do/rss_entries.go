// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RssEntries is the golang structure of table rss_entries for DAO operations like Where/Data.
type RssEntries struct {
	g.Meta    `orm:"table:rss_entries, do:true"`
	Id        any         //
	Guid      any         //
	Url       any         //
	Title     any         //
	Content   any         //
	Published *gtime.Time //
	Author    any         //
}
