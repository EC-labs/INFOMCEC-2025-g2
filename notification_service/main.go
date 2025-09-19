package main

import (
	"log"
	"time"

	"notification_service/config"
	"notification_service/models"
	"notification_service/services"
)

func main() {
	log.Println("Starting Notification Service...")

	// Initialize database connection
	if err := config.ConnectDatabase(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Wait for postgres_service to complete migrations
	log.Println("Waiting for database to be ready...")
	time.Sleep(5 * time.Second)

	// Initialize notification service
	notificationService := services.NewNotificationService()

	// Demonstrate notification service functionality
	demonstrateNotificationService(notificationService)

	log.Println("Notification service is running...")
	
	// Keep the service running
	select {}
}

func demonstrateNotificationService(ns *services.NotificationService) {
	log.Println("\n=== Notification Service Demonstration ===")

	// Create a test user first (this should be done by another service in real scenario)
	db := config.GetDB()
	user := models.User{
		Name:  "Alice Johnson",
		Email: "alice@example.com",
		Age:   28,
	}
	
	// Check if user already exists
	var existingUser models.User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		// User doesn't exist, create new one
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Error creating user: %v", err)
			return
		}
		log.Printf("Created test user: %+v", user)
	} else {
		user = existingUser
		log.Printf("Using existing user: %+v", user)
	}

	// Create a notification template
	template := models.NotificationTemplate{
		Name:     "welcome_email",
		Type:     "email",
		Subject:  "Welcome {{name}}!",
		Template: "Hello {{name}}, welcome to our platform! Your email is {{email}}.",
	}
	
	var existingTemplate models.NotificationTemplate
	if err := db.Where("name = ?", template.Name).First(&existingTemplate).Error; err != nil {
		if err := db.Create(&template).Error; err != nil {
			log.Printf("Error creating template: %v", err)
			return
		}
		log.Printf("Created template: %+v", template)
	} else {
		template = existingTemplate
		log.Printf("Using existing template: %+v", template)
	}

	// 1. Create and send a simple notification
	simpleNotificationReq := models.NotificationRequest{
		UserID:  user.ID,
		Title:   "System Alert",
		Message: "This is a test notification from the notification service!",
		Type:    "push",
	}

	notification1, err := ns.CreateNotification(simpleNotificationReq)
	if err != nil {
		log.Printf("Error creating notification: %v", err)
		return
	}
	log.Printf("Created notification: %+v", notification1)

	// Send the notification
	if err := ns.SendNotification(notification1.ID); err != nil {
		log.Printf("Error sending notification: %v", err)
	}

	// 2. Create notification using template
	templateNotificationReq := models.NotificationRequest{
		UserID:     user.ID,
		Type:       "email",
		TemplateID: &template.ID,
		Variables: map[string]interface{}{
			"name":  user.Name,
			"email": user.Email,
		},
	}

	notification2, err := ns.CreateNotification(templateNotificationReq)
	if err != nil {
		log.Printf("Error creating template notification: %v", err)
		return
	}
	log.Printf("Created template notification: %+v", notification2)

	// Send the template notification
	if err := ns.SendNotification(notification2.ID); err != nil {
		log.Printf("Error sending template notification: %v", err)
	}

	// 3. Create SMS notification
	smsNotificationReq := models.NotificationRequest{
		UserID:  user.ID,
		Title:   "SMS Alert",
		Message: "Your account has been updated successfully!",
		Type:    "sms",
	}

	notification3, err := ns.CreateNotification(smsNotificationReq)
	if err != nil {
		log.Printf("Error creating SMS notification: %v", err)
		return
	}

	if err := ns.SendNotification(notification3.ID); err != nil {
		log.Printf("Error sending SMS notification: %v", err)
	}

	// 4. Retrieve user notifications
	userNotifications, err := ns.GetNotifications(user.ID, 10)
	if err != nil {
		log.Printf("Error getting notifications: %v", err)
		return
	}

	log.Printf("User %d has %d notifications:", user.ID, len(userNotifications))
	for _, notif := range userNotifications {
		log.Printf("  - [%s] %s: %s (Status: %s)", 
			notif.Type, notif.Title, notif.Message, notif.Status)
	}

	// 5. Mark a notification as read
	if len(userNotifications) > 0 {
		if err := ns.MarkAsRead(userNotifications[0].ID); err != nil {
			log.Printf("Error marking notification as read: %v", err)
		} else {
			log.Printf("Marked notification %d as read", userNotifications[0].ID)
		}
	}

	log.Println("\n=== Notification Service Demonstration Complete ===")
}