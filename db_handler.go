package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
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

func addClient(db *sql.DB, client *client) {
	var ip string
	if addr, ok := client.conn.RemoteAddr().(*net.TCPAddr); ok {
		ip = addr.IP.String()
	} else {
		log.Println("Cannot get ip address")
	}
	sqlStatement := `
	INSERT INTO clients_table (ip_address, nickname)
	VALUES ($1, $2)`
	_, err := db.Exec(sqlStatement, ip, client.nickname)
	if err != nil {
		log.Printf("\nError inserting record into database table. %s", err.Error())
	} else {
		log.Println("Successfully added record to database table")
	}

}
