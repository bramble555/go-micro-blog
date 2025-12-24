package model

import "time"

type Message struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	VisitorID string    `gorm:"column:visitor_id;type:varchar(64)" json:"visitor_id"`
	Content   string    `gorm:"column:content;type:varchar(500)" json:"content"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (Message) TableName() string {
	return "messages"
}
