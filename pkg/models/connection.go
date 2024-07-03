package models

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kunalvirwal/go-mvc/pkg/utils"
)

func dataSourceName() string {
	err := godotenv.Load()
	utils.CheckNilErr(err, "Unable to load .env")

	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	database := os.Getenv("DB_DATABASE")

	return fmt.Sprintf("%v:%v@tcp(%v)/%v", username, password, host, database) // Sprintf returns a string instead of printing it to stdout

}

func Connection() (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSourceName()) // db is a connection pool  // here only DSN is validated
	utils.CheckNilErr(err, "unable to open connections")

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second) // we create a context of 5 sec to cancel the connection if nothing returned by pinging the db
	defer cancelfunc()

	err = db.PingContext(ctx) // pings the db for connection verification
	utils.CheckNilErr(err, "DB pinging error")

	fmt.Println("Ayye! DB is connected!")
	return db, err
}
