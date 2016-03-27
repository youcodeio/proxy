package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	id           int
	name         string
	ytid         string
	isTuts       bool
	path         = "/"
	db           *sql.DB
	host         = os.Getenv("HOST_POSTGRE")
	database     = os.Getenv("DATABASE_POSTGRE")
	user         = os.Getenv("USER_POSTGRE")
	password     = os.Getenv("PASSWORD_POSTGRE")
	channelQuery = "SELECT * FROM CHANNELS"
)

// Channel type
type Channel struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Ytid   string `json:"ytid"`
	IsTuts bool   `json:"isTuts"`
}

func init() {
	db = initDB()
}

func getChannels() ([]Channel, error) {
	var channels []Channel
	rows, err := db.Query("select * FROM CHANNELS")
	log.Println("Fetching db")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name, &ytid, &isTuts)
		if err != nil {
			return nil, err
		}
		log.Println(id, name, ytid, isTuts)
		channels = append(channels, Channel{id, name, ytid, isTuts})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func initDB() *sql.DB {
	log.Println("Connecting to postgreSQL")
	dbinfo := fmt.Sprintf("dbname=%s host=%s user=%s password=%s sslmode=disable",
		database, host, user, password)

	DataBase, err := sql.Open("postgres", dbinfo)

	if err != nil {
		log.Fatal(err)
	}

	// Limit conns
	DataBase.SetMaxOpenConns(5)
	return DataBase
}
