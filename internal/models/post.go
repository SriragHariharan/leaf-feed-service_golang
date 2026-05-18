package models

import "time"

//feed_posts table model

type Post struct {
	PostID    string    `json:"post_id" gorm:"primaryKey"`
	MediaURL  string    `json:"media_url" gorm:"not null"`
	Content   string    `json:"content" gorm:"not null"`
	OwnerID   string    `json:"owner_id" gorm:"not null;index:idx_post_owner_created,priority:1"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;index:idx_post_owner_created,priority:2"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`

	Owner *User `gorm:"foreignKey:OwnerID;references:UserID;constraint:-"`
}

func (Post) TableName() string { return "feed_posts" }