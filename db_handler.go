package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func initDatabase() *sql.DB {

	err := godotenv.Load("credentials.env")
	if err != nil {
		log.Printf("\nError when loading the credentials for the database. %s", err.Error())
	}

	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	username := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	db_name := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, db_name)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Printf("\nError when trying to open the database. %s", err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Printf("\nError when trying to ping the database. %s", err.Error())
	}

	return db
}
