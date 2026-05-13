package models

import "time"

//feed_posts table model

type Post struct {
	PostID    string `json:"post_id" gorm:"primaryKey"`
	MediaURL  string `json:"media_url" gorm:"not null"`
	Content   string `json:"content" gorm:"not null"`
	OwnerID   string `json:"owner_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

func (Post) TableName() string { return "feed_posts" }