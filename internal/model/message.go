package model

import "time"

type Message struct {
	ID        int64     `gorm:"primaryKey"`
	VisitorID string    `gorm:"column:visitor_id;type:varchar(64)"`
	Content   string    `gorm:"column:content;type:varchar(500)"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (Message) TableName() string {
	return "messages"
}
