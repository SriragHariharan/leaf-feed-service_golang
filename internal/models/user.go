package models

import "time"

//feed_users table model

type User struct {
	UserID    string `json:"user_id" gorm:"primaryKey"`
	Username  string `json:"username" gorm:"not null"`
	ProfilePic string `json:"profile_pic" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

func (User) TableName() string { return "feed_users" }