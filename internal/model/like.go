package model

import "time"

type Like struct {
	ID        int64     `gorm:"primaryKey" json:"id,string"`
	ArticleID int64     `gorm:"column:article_id;index" json:"article_id,string"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (Like) TableName() string {
	return "likes"
}
