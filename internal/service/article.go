package service

import (
	"errors"

	"go-micro-blog/global"
	"go-micro-blog/internal/model"
)

// CreateArticle 创建文章
func CreateArticle(title string, summary string, content string) (*model.Article, error) {

	// 1️⃣ 基本参数校验
	if title == "" || content == "" {
		return nil, errors.New("title or content cannot be empty")
	}

	// 2️⃣ 使用 Snowflake 生成全局唯一 ID
	article := &model.Article{
		ID:        global.GenID(), // 🔥 核心点
		Title:     title,
		Summary:   summary,
		Content:   content,
		ViewCount: 0,
		Status:    1, // 默认发布
	}

	// 3️⃣ 使用 GORM 写入数据库
	if err := global.DB.Create(article).Error; err != nil {
		return nil, err
	}

	// 4️⃣ 返回创建好的文章（带 ID）
	return article, nil
}
