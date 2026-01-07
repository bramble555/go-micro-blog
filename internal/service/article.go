package service

import (
	"errors"

	"go-micro-blog/global"
	"go-micro-blog/internal/model"

	"go.uber.org/zap"
)

// CreateArticle 创建文章
func CreateArticle(title string, content string) (*model.Article, error) {

	// 1️⃣ 基本参数校验
	if title == "" || content == "" {
		return nil, errors.New("title or content cannot be empty")
	}

	// 2️⃣ 使用 Snowflake 生成全局唯一 ID
	article := &model.Article{
		ID:        global.GenID(), // 🔥 核心点
		Title:     title,
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

func GetArticleList() ([]model.Article, error) {
	var articles []model.Article

	err := global.DB.
		Model(&model.Article{}).
		Where("status = ?", 1).
		Order("created_at DESC").
		Find(&articles).Error

	if err != nil {
		global.Log.Error("查询文章列表失败", zap.Error(err))
		return nil, err
	}
	return articles, nil
}

// GetArticleByID 根据 ID 获取单篇文章详情
func GetArticleByID(id string) (*model.Article, error) {
	var article model.Article
	// id 会自动从 string 转换为数据库匹配的类型
	err := global.DB.Where("id = ?", id).First(&article).Error
	if err != nil {
		global.Log.Error("查询文章详情失败", zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return &article, nil
}

// DeleteArticle 删除文章
func DeleteArticle(id string) error {
	// 1️⃣ 校验 ID 是否存在
	var article model.Article
	if err := global.DB.Where("id = ?", id).First(&article).Error; err != nil {
		return err // 文章不存在，返回错误
	}

	// 2️⃣ 删除文章
	if err := global.DB.Delete(&article).Error; err != nil {
		return err // 删除失败，返回错误
	}

	return nil // 删除成功
}
