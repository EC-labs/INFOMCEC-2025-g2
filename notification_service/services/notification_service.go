package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"notification_service/config"
	"notification_service/models"
	"gorm.io/gorm"
)

// NotificationService handles notification operations
type NotificationService struct {
	db *gorm.DB
}

// NewNotificationService creates a new notification service instance
func NewNotificationService() *NotificationService {
	return &NotificationService{
		db: config.GetDB(),
	}
}

// CreateNotification creates a new notification
func (ns *NotificationService) CreateNotification(req models.NotificationRequest) (*models.Notification, error) {
	// Check if user exists
	var user models.User
	if err := ns.db.First(&user, req.UserID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Prepare metadata
	metadata := ""
	if req.Variables != nil {
		metadataBytes, err := json.Marshal(req.Variables)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata: %w", err)
		}
		metadata = string(metadataBytes)
	}

	// Create notification
	notification := models.Notification{
		UserID:   req.UserID,
		Title:    req.Title,
		Message:  req.Message,
		Type:     req.Type,
		Status:   "pending",
		Metadata: metadata,
	}

	// If template is specified, use it
	if req.TemplateID != nil {
		template, err := ns.GetTemplate(*req.TemplateID)
		if err != nil {
			return nil, fmt.Errorf("failed to get template: %w", err)
		}
		
		// Process template with variables
		processedMessage, err := ns.ProcessTemplate(template.Template, req.Variables)
		if err != nil {
			return nil, fmt.Errorf("failed to process template: %w", err)
		}
		
		notification.Message = processedMessage
		if template.Subject != "" {
			notification.Title = template.Subject
		}
	}

	if err := ns.db.Create(&notification).Error; err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return &notification, nil
}

// SendNotification simulates sending a notification
func (ns *NotificationService) SendNotification(notificationID uint) error {
	var notification models.Notification
	if err := ns.db.First(&notification, notificationID).Error; err != nil {
		return fmt.Errorf("notification not found: %w", err)
	}

	// Simulate sending based on type
	switch notification.Type {
	case "email":
		log.Printf("📧 Sending email to user %d: %s", notification.UserID, notification.Title)
	case "sms":
		log.Printf("📱 Sending SMS to user %d: %s", notification.UserID, notification.Message)
	case "push":
		log.Printf("🔔 Sending push notification to user %d: %s", notification.UserID, notification.Title)
	default:
		return fmt.Errorf("unsupported notification type: %s", notification.Type)
	}

	// Update notification status
	now := time.Now()
	notification.Status = "sent"
	notification.SentAt = &now

	if err := ns.db.Save(&notification).Error; err != nil {
		return fmt.Errorf("failed to update notification status: %w", err)
	}

	log.Printf("✅ Notification %d sent successfully", notificationID)
	return nil
}

// GetNotifications retrieves notifications for a user
func (ns *NotificationService) GetNotifications(userID uint, limit int) ([]models.Notification, error) {
	var notifications []models.Notification
	query := ns.db.Where("user_id = ?", userID).Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&notifications).Error; err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	return notifications, nil
}

// GetTemplate retrieves a notification template
func (ns *NotificationService) GetTemplate(templateID uint) (*models.NotificationTemplate, error) {
	var template models.NotificationTemplate
	if err := ns.db.First(&template, templateID).Error; err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}
	return &template, nil
}

// ProcessTemplate processes a template with variables (simple string replacement)
func (ns *NotificationService) ProcessTemplate(template string, variables map[string]interface{}) (string, error) {
	result := template
	
	// Simple template processing - replace {{variable}} with values
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		valueStr := fmt.Sprintf("%v", value)
		result = strings.ReplaceAll(result, placeholder, valueStr)
	}
	
	return result, nil
}

// MarkAsRead marks a notification as read
func (ns *NotificationService) MarkAsRead(notificationID uint) error {
	result := ns.db.Model(&models.Notification{}).
		Where("id = ?", notificationID).
		Update("status", "read")
	
	if result.Error != nil {
		return fmt.Errorf("failed to mark notification as read: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}
	
	return nil
}