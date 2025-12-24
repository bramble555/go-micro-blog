package model

import "time"

type Article struct {
	ID        int64     `json:"id,string" gorm:"primaryKey"`
	Title     string    `gorm:"column:title;type:varchar(200);not null" json:"title"`
	Summary   string    `gorm:"column:summary;type:varchar(500)" json:"summary"`
	Content   string    `gorm:"column:content;type:longtext;not null" json:"content"`
	ViewCount int64     `gorm:"column:view_count;default:0" json:"view_count"`
	Status    int8      `gorm:"column:status;default:1" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Article) TableName() string {
	return "articles"
}
