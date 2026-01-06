package model

import "time"

type Comment struct {
	ID        int64     `gorm:"primaryKey" json:"id,string"`
	ArticleID int64     `gorm:"column:article_id;index" json:"article_id,string"`
	Nickname  string    `gorm:"column:nickname;type:varchar(100);default:'游客'" json:"nickname"`
	Content   string    `gorm:"column:content;type:varchar(500)" json:"content"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (Comment) TableName() string {
	return "comments"
}
