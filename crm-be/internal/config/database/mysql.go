package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbmysql *gorm.DB

func DatabaseMysql() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("MYSQL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	// db.AutoMigrate(&User{}, &Product{})

	dbmysql = db

	return db
}

func GetDB() *gorm.DB {
	if dbmysql == nil {
		log.Println("Reload Connection", dbmysql)
		DatabaseMysql()
	}
	return dbmysql
}
