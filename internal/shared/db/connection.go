package db

import (
	"fmt"
	"log"
	"os"
	"socialmedia-backend/internal/shared/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func createDB() {
	err := db.AutoMigrate(
		&models.User{},
		&models.Conversation{},
		&models.Message{},
		&models.UnreadMessage{},
	)
	if err != nil {
		log.Fatal("migrations failed: ", err)
	}
	fmt.Println("migrations completed successfully")
}

func InitDB() {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("socialmedia_db connection failed")
	}

	db = database

	createDB()

	fmt.Println("socialmedia_db connection successful")
}

func GetDB() *gorm.DB {
	return db
}
