package models

import (
	"time"

	"gorm.io/gorm"
)

// These models mirror the ones in postgres_service for database interaction
// The postgres_service manages migrations and schema changes

// User represents a user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Age       int            `json:"age"`
	Notifications []Notification `json:"notifications" gorm:"foreignKey:UserID"`
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

// NotificationRequest represents a request to send a notification
type NotificationRequest struct {
	UserID     uint                   `json:"user_id" validate:"required"`
	Title      string                 `json:"title" validate:"required"`
	Message    string                 `json:"message" validate:"required"`
	Type       string                 `json:"type" validate:"required,oneof=email sms push"`
	TemplateID *uint                  `json:"template_id,omitempty"`
	Variables  map[string]interface{} `json:"variables,omitempty"`
}