-- RSS条目表索引优化
-- 适用于 6万+ 数据量

-- 1. 发布时间索引（排序和日期范围查询最关键）
CREATE INDEX idx_published ON rss_entries(published DESC);

-- 2. 标题索引（用于模糊搜索）
CREATE INDEX idx_title ON rss_entries(title(100));

-- 3. 联合索引（标题搜索 + 时间排序）
-- 注意：如果已创建 idx_title，可以先删除
-- DROP INDEX idx_title ON rss_entries;
-- CREATE INDEX idx_title_published ON rss_entries(title(100), published DESC);

-- 4. 可选：全文索引（MySQL 5.7+，支持中文需配置ngram分词器）
-- ALTER TABLE rss_entries ADD FULLTEXT INDEX ft_title(title) WITH PARSER ngram;

-- 查看索引
SHOW INDEX FROM rss_entries;

-- 检查查询性能
EXPLAIN SELECT * FROM rss_entries WHERE title LIKE '%关键词%' ORDER BY published DESC LIMIT 20;
