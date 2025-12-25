package model

import "time"

type Comment struct {
	ID        int64     `gorm:"primaryKey" json:"id,string"`
	ArticleID int64     `gorm:"column:article_id;index" json:"article_id,string"`
	VisitorID string    `gorm:"column:visitor_id;type:varchar(64)" json:"visitor_id"`
	Content   string    `gorm:"column:content;type:varchar(500)" json:"content"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (Comment) TableName() string {
	return "comments"
}
