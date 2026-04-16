package initdb

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

var CreateTableSQL = "CREATE TABLE IF NOT EXISTS `rss_entries` (" +
	"`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT," +
	"`guid` VARCHAR(255) NOT NULL COMMENT '唯一标识符'," +
	"`url` VARCHAR(512) NOT NULL COMMENT '文章链接'," +
	"`title` VARCHAR(255) NOT NULL COMMENT '文章标题'," +
	"`content` LONGTEXT COMMENT '文章内容'," +
	"`published` DATETIME COMMENT '发布时间'," +
	"`author` VARCHAR(100) DEFAULT NULL COMMENT '作者'," +
	"`created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'," +
	"`updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'," +
	"PRIMARY KEY (`id`)," +
	"UNIQUE KEY `uk_guid` (`guid`)," +
	"KEY `idx_published` (`published` DESC)," +
	"KEY `idx_title` (`title`)" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='RSS条目表'"

func EnsureTables(ctx context.Context) error {
	_, err := g.DB().Exec(ctx, CreateTableSQL)
	return err
}
