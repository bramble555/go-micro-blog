package service

import (
	"go-micro-blog/global"
	"go-micro-blog/internal/model"
)

func GetLikesByArticleID(articleID int64) (int64, error) {
	var count int64
	err := global.DB.Model(&model.Like{}).Where("article_id = ?", articleID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
func CreateLike(articleID int64) error {
	like := model.Like{
		ID:        global.GenID(),
		ArticleID: articleID,
	}
	err := global.DB.Create(&like).Error
	if err != nil {
		return err
	}
	return nil
}
