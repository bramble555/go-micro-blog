package model

import "time"

type Like struct {
	ID        int64     `gorm:"primaryKey"`
	ArticleID int64     `gorm:"column:article_id;index"`
	VisitorID string    `gorm:"column:visitor_id;type:varchar(64)"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (Like) TableName() string {
	return "likes"
}
