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

type db_client struct {
	ip       string
	nickname string
}

type db_handler struct {
	db         *sql.DB
	db_clients []db_client
}

func (dbh *db_handler) initDatabase() {

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

	dbh.db = db
}

func (dbh *db_handler) addClient(client db_client) {
	sqlStatement := `
	INSERT INTO clients_table (ip_address, nickname)
	VALUES ($1, $2)`
	_, err := dbh.db.Exec(sqlStatement, client.ip, client.nickname)
	if err != nil {
		log.Printf("\nError inserting record into database table. %s", err.Error())
	} else {
		log.Println("Successfully added record to database table")
	}

}

func (dbh *db_handler) updateClientRecord(client db_client) {
	sqlStatement := `
	UPDATE clients_table
	SET nickname = $1
	WHERE ip_address = $2;`

	_, err := dbh.db.Exec(sqlStatement, client.nickname, client.ip)
	if err != nil {
		log.Printf("Error updating client record. %s", err.Error())
	}
}

func (dbh *db_handler) pullClients() {
	sqlStatement := `SELECT * FROM clients_table`
	var res []db_client
	rows, err := dbh.db.Query(sqlStatement)
	if err != nil {
		log.Printf("Error getting records from table. %s", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var ip string
		var nick string
		err = rows.Scan(&ip, &nick)
		if err != nil {
			panic(err)
		}
		res = append(res, db_client{
			ip:       ip,
			nickname: nick,
		})
	}

	dbh.db_clients = res
}
