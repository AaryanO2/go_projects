package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func Connect() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	var d *gorm.DB
	var err error
	// Retry loop for MySQL connection (max 10 attempts)
	for i := 0; i < 10; i++ {
		d, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			db = d
			return
		}
		log.Printf("MySQL not ready, retrying in 2 seconds (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	// If max retries reached, panic with the error.
	panic(fmt.Errorf("failed to connect to database: %w", err))
}

func GetDB() *gorm.DB {
	return db
}
