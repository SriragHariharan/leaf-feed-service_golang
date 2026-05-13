package models

import "time"

type Feed struct {
	FeedID    string    `json:"feed_id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null"`
	PostID    string    `json:"post_id" gorm:"not null"`
	IsLiked   bool      `json:"is_liked" gorm:"not null"`
	IsCommented bool      `json:"is_commented" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`

	User *User `gorm:"foreignKey:UserID"`
	Post *Post `gorm:"foreignKey:PostID"`
}

func (Feed) TableName() string { return "feed_feeds" }