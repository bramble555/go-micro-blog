-- 给 comments 表添加 nickname 字段（如果表已存在）
ALTER TABLE comments ADD COLUMN nickname VARCHAR(100) DEFAULT '游客' AFTER article_id;

-- 移除 visitor_id 字段（如果存在）
ALTER TABLE comments DROP COLUMN IF EXISTS visitor_id;
