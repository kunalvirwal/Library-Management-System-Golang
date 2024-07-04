package models

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kunalvirwal/go-mvc/pkg/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func dataSourceName() string {
	err := godotenv.Load()
	utils.CheckNilErr(err, "Unable to load .env")

	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	database := os.Getenv("DB_DATABASE")

	return fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, database) // Sprintf returns a string instead of printing it to stdout

}

func Connection() (*gorm.DB, error) {

	db, err := gorm.Open(mysql.Open(dataSourceName()), &gorm.Config{}) // db is a connection pool  // here only DSN is validated
	utils.CheckNilErr(err, "Unable to connect to DB")

	db.AutoMigrate(&BOOKS{})
	db.AutoMigrate(&USER{})
	db.AutoMigrate(&BORROWING_HISTORY{})
	db.AutoMigrate(&PENDING_REQUESTS{})

	sqldb, err := db.DB()
	utils.CheckNilErr(err, "Unable to get SQL DB instance")

	sqldb.SetMaxOpenConns(20)
	sqldb.SetMaxIdleConns(20)
	sqldb.SetConnMaxLifetime(time.Minute * 5)

	return db, err
}
