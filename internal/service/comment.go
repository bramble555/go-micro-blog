package service

import (
	"go-micro-blog/global"
	"go-micro-blog/internal/model"
)

// 获取文章评论列表
func GetCommentsByArticleID(articleID int64) ([]model.Comment, error) {
	var comments []model.Comment

	err := global.DB.
		Where("article_id = ?", articleID).
		Order("created_at DESC").
		Find(&comments).Error

	if err != nil {
		return nil, err
	}

	return comments, nil
}

// 创建评论
func CreateComment(articleID int64, nickname, content string) error {
	comment := model.Comment{
		ID:        global.GenID(),
		ArticleID: articleID,
		// Nickname:  nickname,
		Content: content,
	}

	return global.DB.Create(&comment).Error
}

// 管理员删除评论
func DeleteComment(commentID int64) error {
	return global.DB.Delete(&model.Comment{}, commentID).Error
}
