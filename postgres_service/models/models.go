package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Age       int            `json:"age"`
	Posts     []Post         `json:"posts" gorm:"foreignKey:UserID"`
	Notifications []Notification `json:"notifications" gorm:"foreignKey:UserID"`
}

// Post represents a blog post
type Post struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Title     string         `json:"title" gorm:"not null"`
	Content   string         `json:"content" gorm:"type:text"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Tags      []Tag          `json:"tags" gorm:"many2many:post_tags;"`
}

// Tag represents a tag for posts
type Tag struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Name      string         `json:"name" gorm:"uniqueIndex;not null"`
	Posts     []Post         `json:"posts" gorm:"many2many:post_tags;"`
}

// Notification represents a notification sent to users
type Notification struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Title     string         `json:"title" gorm:"not null"`
	Message   string         `json:"message" gorm:"type:text;not null"`
	Type      string         `json:"type" gorm:"not null"` // email, sms, push, etc.
	Status    string         `json:"status" gorm:"default:'pending'"` // pending, sent, failed
	SentAt    *time.Time     `json:"sent_at"`
	Metadata  string         `json:"metadata" gorm:"type:jsonb"` // Additional data as JSON
}

// NotificationTemplate represents reusable notification templates
type NotificationTemplate struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Name      string         `json:"name" gorm:"uniqueIndex;not null"`
	Type      string         `json:"type" gorm:"not null"`
	Subject   string         `json:"subject"`
	Template  string         `json:"template" gorm:"type:text;not null"`
	Variables string         `json:"variables" gorm:"type:jsonb"` // Template variables as JSON
}

// GetAllModels returns all models for migration
func GetAllModels() []interface{} {
	return []interface{}{
		&User{},
		&Post{},
		&Tag{},
		&Notification{},
		&NotificationTemplate{},
	}
}