package rss_sync

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cicbyte/itfeeds/internal/dao"
	"github.com/cicbyte/itfeeds/internal/model/do"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/mmcdole/gofeed"
)

var httpClient = &http.Client{Timeout: 15 * time.Second}

// StartRSSSync 启动RSS同步定时任务
func StartRSSSync(ctx context.Context) error {
	enabled := g.Cfg().MustGet(ctx, "rss.enabled", false).Bool()
	if !enabled {
		g.Log().Info(ctx, "[RSS] RSS同步已禁用")
		return nil
	}

	cronExprs := g.Cfg().MustGet(ctx, "rss.crons").Strings()
	if len(cronExprs) == 0 {
		singleCron := g.Cfg().MustGet(ctx, "rss.cron", "").String()
		if singleCron != "" {
			cronExprs = []string{singleCron}
		}
	}

	feeds := g.Cfg().MustGet(ctx, "rss.feeds").Strings()

	if len(feeds) == 0 {
		g.Log().Warning(ctx, "[RSS] 未配置RSS源")
		return nil
	}

	if len(cronExprs) == 0 {
		g.Log().Warning(ctx, "[RSS] 未配置定时规则")
		return nil
	}

	g.Log().Infof(ctx, "[RSS] 启动定时同步，cron规则: %v, 源数量: %d", cronExprs, len(feeds))

	for i, cronExpr := range cronExprs {
		taskName := fmt.Sprintf("rss_sync_%d", i)
		_, err := gcron.Add(ctx, cronExpr, func(ctx context.Context) {
			syncAllFeeds(ctx, feeds)
		}, taskName)
		if err != nil {
			g.Log().Errorf(ctx, "[RSS] 注册定时任务失败 [%s]: %v", cronExpr, err)
			return err
		}
		g.Log().Infof(ctx, "[RSS] 定时任务已注册: %s -> %s", taskName, cronExpr)
	}

	go syncAllFeeds(context.Background(), feeds)

	return nil
}

// SyncAllFeeds 同步所有RSS源
func SyncAllFeeds(ctx context.Context) int {
	feeds := g.Cfg().MustGet(ctx, "rss.feeds").Strings()
	if len(feeds) == 0 {
		return 0
	}

	totalAdded := 0
	for _, feedURL := range feeds {
		added, err := syncFeed(ctx, feedURL)
		if err != nil {
			g.Log().Errorf(ctx, "[RSS] 同步失败 %s: %v", feedURL, err)
			continue
		}
		totalAdded += added
	}

	if totalAdded > 0 {
		sendPush(ctx, totalAdded)
	}
	return totalAdded
}

// syncAllFeeds 同步所有RSS源（定时任务内部使用）
func syncAllFeeds(ctx context.Context, feeds []string) {
	startTime := time.Now()
	g.Log().Info(ctx, "[RSS] 开始同步...")

	totalAdded := 0
	for _, feedURL := range feeds {
		added, err := syncFeed(ctx, feedURL)
		if err != nil {
			g.Log().Errorf(ctx, "[RSS] 同步失败 %s: %v", feedURL, err)
			continue
		}
		totalAdded += added
	}

	elapsed := time.Since(startTime)
	g.Log().Infof(ctx, "[RSS] 同步完成，新增 %d 条，耗时 %v", totalAdded, elapsed)

	if totalAdded > 0 {
		sendPush(ctx, totalAdded)
	}
}

// syncFeed 同步单个RSS源，利用唯一索引批量插入
func syncFeed(ctx context.Context, feedURL string) (int, error) {
	fp := gofeed.NewParser()
	fp.Client = httpClient
	feed, err := fp.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return 0, fmt.Errorf("解析RSS失败: %w", err)
	}

	// 收集所有待插入条目
	var items []do.RssEntries
	for _, item := range feed.Items {
		guid := item.GUID
		if guid == "" {
			guid = item.Link
		}
		if guid == "" {
			continue
		}

		var published *gtime.Time
		if item.PublishedParsed != nil {
			published = gtime.New(item.PublishedParsed.Local())
		}

		content := item.Content
		if content == "" {
			content = item.Description
		}

		items = append(items, do.RssEntries{
			Guid:      guid,
			Url:       item.Link,
			Title:     item.Title,
			Content:   content,
			Published: published,
			Author:    getAuthor(item),
		})
	}

	if len(items) == 0 {
		return 0, nil
	}

	// 批量插入，利用 uk_guid 唯一索引忽略重复
	result, err := dao.RssEntries.Ctx(ctx).InsertIgnore(items)
	if err != nil {
		return 0, fmt.Errorf("批量插入失败: %w", err)
	}

	added := 0
	if result != nil {
		rowsAffected, _ := result.RowsAffected()
		added = int(rowsAffected)
	}

	g.Log().Infof(ctx, "[RSS] %s: 处理 %d 条，新增 %d 条", feedURL, len(items), added)
	return added, nil
}

// getAuthor 获取作者
func getAuthor(item *gofeed.Item) string {
	if item.Author != nil {
		return item.Author.Name
	}
	return ""
}

// sendPush 发送推送通知
func sendPush(ctx context.Context, count int) {
	barkKey := g.Cfg().MustGet(ctx, "rss.barkPush", "").String()
	if barkKey == "" {
		return
	}

	title := "ITFeeds 同步完成"
	body := fmt.Sprintf("新增 %d 篇文章", count)
	url := fmt.Sprintf("https://api.day.app/%s/%s/%s", barkKey, title, body)

	resp, err := httpClient.Get(url)
	if err != nil {
		g.Log().Errorf(ctx, "[RSS] 推送失败: %v", err)
		return
	}
	defer resp.Body.Close()

	g.Log().Info(ctx, "[RSS] 推送通知已发送")
}
