// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RssEntriesDao is the data access object for the table rss_entries.
type RssEntriesDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  RssEntriesColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// RssEntriesColumns defines and stores column names for the table rss_entries.
type RssEntriesColumns struct {
	Id        string //
	Guid      string //
	Url       string //
	Title     string //
	Content   string //
	Published string //
	Author    string //
}

// rssEntriesColumns holds the columns for the table rss_entries.
var rssEntriesColumns = RssEntriesColumns{
	Id:        "id",
	Guid:      "guid",
	Url:       "url",
	Title:     "title",
	Content:   "content",
	Published: "published",
	Author:    "author",
}

// NewRssEntriesDao creates and returns a new DAO object for table data access.
func NewRssEntriesDao(handlers ...gdb.ModelHandler) *RssEntriesDao {
	return &RssEntriesDao{
		group:    "default",
		table:    "rss_entries",
		columns:  rssEntriesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *RssEntriesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *RssEntriesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *RssEntriesDao) Columns() RssEntriesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *RssEntriesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *RssEntriesDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *RssEntriesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
