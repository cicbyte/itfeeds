package mcp

import (
	"context"
	"fmt"
	"runtime"
	"time"

	api "github.com/cicbyte/itfeeds/api/v1/rss_entries"
	"github.com/cicbyte/itfeeds/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerTools 注册所有工具到 MCP 服务器
func registerTools(s *server.MCPServer) {
	// 1. 获取服务器信息
	s.AddTool(mcp.NewTool("get_server_info",
		mcp.WithDescription("获取 MCP 服务器的运行信息"),
	), handleGetServerInfo)

	// 2. 获取新闻列表
	s.AddTool(mcp.NewTool("get_news_list",
		mcp.WithDescription("获取ITFeeds新闻列表，支持分页、标题搜索和时间范围筛选"),
		mcp.WithNumber("page_num",
			mcp.Description("页码，默认1"),
		),
		mcp.WithNumber("page_size",
			mcp.Description("每页数量，默认20"),
		),
		mcp.WithString("title",
			mcp.Description("标题关键词搜索"),
		),
		mcp.WithString("start_date",
			mcp.Description("开始日期，格式：YYYY-MM-DD"),
		),
		mcp.WithString("end_date",
			mcp.Description("结束日期，格式：YYYY-MM-DD"),
		),
	), handleGetNewsList)

	// 3. 获取新闻详情
	s.AddTool(mcp.NewTool("get_news_detail",
		mcp.WithDescription("根据ID获取新闻详情"),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description("新闻ID"),
		),
	), handleGetNewsDetail)

	// 4. 搜索新闻
	s.AddTool(mcp.NewTool("search_news",
		mcp.WithDescription("搜索新闻标题，返回匹配的新闻列表"),
		mcp.WithString("keyword",
			mcp.Required(),
			mcp.Description("搜索关键词"),
		),
		mcp.WithNumber("limit",
			mcp.Description("返回数量限制，默认10"),
		),
	), handleSearchNews)

	// 5. 获取统计信息
	s.AddTool(mcp.NewTool("get_statistics",
		mcp.WithDescription("获取新闻数据统计信息"),
	), handleGetStatistics)
}

// handleGetServerInfo 处理获取服务器信息
func handleGetServerInfo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	info := fmt.Sprintf(
		"ITFeeds MCP Server:\n"+
			"  Version: 1.0.0\n"+
			"  Go Version: %s\n"+
			"  OS: %s\n"+
			"  Architecture: %s\n"+
			"  CPU Cores: %d\n"+
			"  Current Time: %s",
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
		time.Now().Format("2006-01-02 15:04:05"),
	)

	return mcp.NewToolResultText(info), nil
}

// handleGetNewsList 处理获取新闻列表
func handleGetNewsList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	// 默认值
	pageNum := 1
	pageSize := 20

	if pn, ok := args["page_num"].(float64); ok && pn > 0 {
		pageNum = int(pn)
	}
	if ps, ok := args["page_size"].(float64); ok && ps > 0 {
		pageSize = int(ps)
	}

	// 构建请求
	req := &api.RssEntriesListReq{}
	req.PageReq.PageNum = pageNum
	req.PageReq.PageSize = pageSize

	if title, ok := args["title"].(string); ok {
		req.Title = title
	}
	if startDate, ok := args["start_date"].(string); ok {
		req.StartDate = startDate
	}
	if endDate, ok := args["end_date"].(string); ok {
		req.EndDate = endDate
	}

	// 调用服务层
	total, list, err := service.RssEntries().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("获取新闻列表失败: %v", err)
	}

	// 格式化输出
	var result string
	result += fmt.Sprintf("共找到 %v 条新闻，当前第 %d 页，每页 %d 条\n\n", total, pageNum, pageSize)

	for i, item := range list {
		published := ""
		if item.Published != nil {
			published = item.Published.Format("2006-01-02 15:04")
		}
		result += fmt.Sprintf("%d. [%d] %s\n", (pageNum-1)*pageSize+i+1, item.Id, item.Title)
		result += fmt.Sprintf("   作者: %s | 发布时间: %s\n", item.Author, published)
		if i < len(list)-1 {
			result += "\n"
		}
	}

	return mcp.NewToolResultText(result), nil
}

// handleGetNewsDetail 处理获取新闻详情
func handleGetNewsDetail(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	id, ok := args["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("id 参数必须是数字")
	}

	// 调用服务层
	detail, err := service.RssEntries().GetById(ctx, int(id))
	if err != nil {
		return nil, fmt.Errorf("获取新闻详情失败: %v", err)
	}

	// 格式化输出
	published := ""
	if detail.Published != nil {
		published = detail.Published.Format("2006-01-02 15:04:05")
	}

	result := fmt.Sprintf(
		"【新闻详情】\n"+
			"ID: %d\n"+
			"标题: %s\n"+
			"作者: %s\n"+
			"发布时间: %s\n"+
			"链接: %s\n\n"+
			"内容:\n%s",
		detail.Id,
		detail.Title,
		detail.Author,
		published,
		detail.Url,
		detail.Content,
	)

	return mcp.NewToolResultText(result), nil
}

// handleSearchNews 处理搜索新闻
func handleSearchNews(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	keyword, ok := args["keyword"].(string)
	if !ok || keyword == "" {
		return nil, fmt.Errorf("keyword 参数不能为空")
	}

	limit := 10
	if l, ok := args["limit"].(float64); ok && l > 0 {
		limit = int(l)
	}

	// 调用服务层
	req := &api.RssEntriesListReq{}
	req.PageReq.PageNum = 1
	req.PageReq.PageSize = limit
	req.Title = keyword

	total, list, err := service.RssEntries().List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("搜索新闻失败: %v", err)
	}

	// 格式化输出
	result := fmt.Sprintf("搜索 \"%s\" 找到 %v 条结果:\n\n", keyword, total)

	for i, item := range list {
		published := ""
		if item.Published != nil {
			published = item.Published.Format("2006-01-02")
		}
		result += fmt.Sprintf("%d. [%d] %s\n", i+1, item.Id, item.Title)
		result += fmt.Sprintf("   %s | %s\n", item.Author, published)
		if i < len(list)-1 {
			result += "\n"
		}
	}

	return mcp.NewToolResultText(result), nil
}

// handleGetStatistics 处理获取统计信息
func handleGetStatistics(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 一次查询获取全部统计
	sql := `SELECT
		COUNT(*) AS total,
		COALESCE(SUM(CASE WHEN published >= CURDATE() THEN 1 ELSE 0 END), 0) AS today,
		COALESCE(SUM(CASE WHEN published >= DATE_SUB(CURDATE(), INTERVAL 6 DAY) THEN 1 ELSE 0 END), 0) AS week
	FROM rss_entries`

	type stats struct {
		Total int `json:"total"`
		Today int `json:"today"`
		Week  int `json:"week"`
	}

	var s stats
	err := g.DB().Raw(sql).Scan(ctx, &s)
	if err != nil {
		return nil, fmt.Errorf("获取统计信息失败: %v", err)
	}

	result := fmt.Sprintf(
		"【ITFeeds 新闻统计】\n总新闻数: %d\n今日新增: %d\n最近7天: %d\n统计时间: %s",
		s.Total, s.Today, s.Week,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	return mcp.NewToolResultText(result), nil
}