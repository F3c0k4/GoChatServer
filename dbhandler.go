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

// dbClient stores the information of a client
// in a way that mirrors a database record
type dbClient struct {
	ip       string
	nickname string
}

// dbHandler contains information which is useful
// in contacting the database and keeping track of the clients
type dbHandler struct {
	db        *sql.DB
	dbClients []dbClient
}

// initDatabase loads the credentials of the database from the credentials.env file
// and attempts to connect to the database
func (dbh *dbHandler) initDatabase() error {

	err := godotenv.Load("credentials.env")
	if err != nil {
		return fmt.Errorf("\nError when loading the credentials for the database. %w", err)
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
		return fmt.Errorf("\nError when trying to open the database. %w", err)
	}

	dbh.db = db
	return nil
}

// getClient looks up a client based on the ip parameter it receives
// and returns a pointer to the db_client
func (dbh *dbHandler) getClient(ip string) *dbClient {
	for _, c := range dbh.dbClients {
		if c.ip == ip {
			return &c
		}
	}

	return nil
}

// addClient adds a new client record to the database
func (dbh *dbHandler) addClient(client dbClient) {
	sqlStatement := `
	INSERT INTO clients_table (ip_address, nickname)
	VALUES ($1, $2)`
	_, err := dbh.db.Exec(sqlStatement, client.ip, client.nickname)
	if err != nil {
		log.Printf("\nError inserting record into database table. %s", err.Error())
	} else {
		log.Println("Successfully added record to database table")
	}

	dbh.dbClients = append(dbh.dbClients, client)
}

// updateClientRecord updates the database with the data
// of the client object it receives
func (dbh *dbHandler) updateClientRecord(client dbClient) {
	sqlStatement := `
	UPDATE clients_table
	SET nickname = $1
	WHERE ip_address = $2;`

	_, err := dbh.db.Exec(sqlStatement, client.nickname, client.ip)
	if err != nil {
		log.Printf("Error updating client record. %s", err.Error())
	}
}

// pullClients loads data from the database into the
// dbHandlers dbClients slice
func (dbh *dbHandler) pullClients() {
	sqlStatement := `SELECT * FROM clients_table`
	var res []dbClient
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
		res = append(res, dbClient{
			ip:       ip,
			nickname: nick,
		})
	}

	dbh.dbClients = res
}
