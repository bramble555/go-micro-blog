package model

import "time"

type Article struct {
	ID        int64     `gorm:"primaryKey;column:id"`
	Title     string    `gorm:"column:title;type:varchar(200);not null"`
	Summary   string    `gorm:"column:summary;type:varchar(500)"`
	Content   string    `gorm:"column:content;type:longtext;not null"`
	ViewCount int64     `gorm:"column:view_count;default:0"`
	Status    int8      `gorm:"column:status;default:1"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Article) TableName() string {
	return "articles"
}
