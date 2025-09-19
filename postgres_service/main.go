package main

import (
	"flag"
	"log"
	"postgres_service/config"
	"postgres_service/migrations"
)

func main() {
	// Parse command line flags
	migrate := flag.Bool("migrate", false, "Run database migrations")
	demo := flag.Bool("demo", false, "Run CRUD demonstration")
	flag.Parse()

	// Initialize database connection
	if err := config.ConnectDatabase(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations if requested
	if *migrate {
		if err := migrations.RunMigrations(); err != nil {
			log.Fatal("Failed to run migrations:", err)
		}
		log.Println("Migrations completed successfully!")
		return
	}

	// Run demo if requested
	if *demo {
		demonstrateCRUD()
		return
	}

	// Default behavior: run migrations and keep service running
	if err := migrations.RunMigrations(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Database management service started. Migrations completed.")
	log.Println("Database is ready for other services to connect.")
	
	// Keep the service running (for Docker health checks)
	select {}
}

func demonstrateCRUD() {
	db := config.GetDB()

	fmt.Println("\n=== GORM CRUD Demonstration ===")

	// Create a user
	user := models.User{
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Age:   30,
	}

	result := db.Create(&user)
	if result.Error != nil {
		log.Printf("Error creating user: %v", result.Error)
		return
	}
	fmt.Printf("Created user: %+v\n", user)

	// Create a post for the user
	post := models.Post{
		Title:   "My First Post",
		Content: "This is the content of my first post using GORM!",
		UserID:  user.ID,
	}

	if err := db.Create(&post).Error; err != nil {
		log.Printf("Error creating post: %v", err)
		return
	}
	fmt.Printf("Created post: %+v\n", post)

	// Create tags
	tag1 := models.Tag{Name: "technology"}
	tag2 := models.Tag{Name: "golang"}

	db.Create(&tag1)
	db.Create(&tag2)

	// Associate tags with the post
	db.Model(&post).Association("Tags").Append([]models.Tag{tag1, tag2})

	// Read - Find user by ID with posts
	var foundUser models.User
	db.Preload("Posts").First(&foundUser, user.ID)
	fmt.Printf("Found user with posts: %+v\n", foundUser)

	// Read - Find post with user and tags
	var foundPost models.Post
	db.Preload("User").Preload("Tags").First(&foundPost, post.ID)
	fmt.Printf("Found post with user and tags: %+v\n", foundPost)

	// Update - Update user's age
	db.Model(&user).Update("age", 31)
	fmt.Printf("Updated user age to 31\n")

	// Read all users
	var users []models.User
	db.Find(&users)
	fmt.Printf("All users: %+v\n", users)

	// Count records
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	fmt.Printf("Total users in database: %d\n", userCount)

	// Custom query example
	var adultUsers []models.User
	db.Where("age >= ?", 18).Find(&adultUsers)
	fmt.Printf("Adult users: %+v\n", adultUsers)

	// Soft delete example (won't actually delete from database)
	db.Delete(&post)
	fmt.Printf("Soft deleted post with ID: %d\n", post.ID)

	// Verify soft delete - this won't find the deleted post
	var activePosts []models.Post
	db.Find(&activePosts)
	fmt.Printf("Active posts count: %d\n", len(activePosts))

	// Find deleted posts (including soft deleted)
	var allPosts []models.Post
	db.Unscoped().Find(&allPosts)
	fmt.Printf("All posts (including deleted) count: %d\n", len(allPosts))

	fmt.Println("\n=== CRUD Demonstration Complete ===")
}