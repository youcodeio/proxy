package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var (
	// Cache version of ChannelsLists. Upgrade every hour, see init()
	channelsListsCached []Channel
	id                  int
	name                string
	ytid                string
	isTuts              bool
	path                = "/"
	host                = os.Getenv("POSTGRESQL_ADDON_HOST")
	database            = os.Getenv("POSTGRESQL_ADDON_DB")
	user                = os.Getenv("POSTGRESQL_ADDON_USER")
	password            = os.Getenv("POSTGRESQL_ADDON_PASSWORD")
)

type YouCodeDB struct {
	*sql.DB
}

// Channel type
type Channel struct {
	Name   string `json:"name"`
	Ytid   string `json:"ytid"`
	IsTuts bool   `json:"isTuts"`
}

type AppDatabase interface {
	getChannelsFromDB() []Channel
}

//AddChannels provide a simple way for admin to add channels
func (db *YouCodeDB) AddChannel(name string, ytid string, isTuts bool) bool {

	log.Println("Inserting", name, ytid, isTuts)

	stmt, err := db.Prepare("INSERT INTO CHANNELS (NAME, YTID, ISTUTS) VALUES (?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(name, ytid, isTuts)
	if err != nil {
		log.Fatal(err)
		return false
	}
	log.Println(res)

	return true
}

func (db *YouCodeDB) GetChannels() []Channel {
	return channelsListsCached
}

func (db *YouCodeDB) getChannelsFromDB() []Channel {
	var channels []Channel
	rows, err := db.Query("select * FROM CHANNELS")
	log.Println("Fetching db...")
	if err != nil {
		log.Fatal("Error fetching Database:", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name, &ytid, &isTuts)
		if err != nil {
			log.Fatal("Error fetching Database:", err)
		}
		channels = append(channels, Channel{name, ytid, isTuts})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal("Error fetching Database:", err)
	}
	log.Println("Fetching db done, found", len(channels), "channels")
	return channels
}

func InitDB() *YouCodeDB {
	log.Println("Connecting to postgreSQL")
	dbinfo := fmt.Sprintf("dbname=%s host=%s user=%s password=%s sslmode=disable",
		database, host, user, password)

	DataBase, err := sql.Open("postgres", dbinfo)

	if err != nil {
		log.Fatal(err)
	}

	// Limit conns
	DataBase.SetMaxOpenConns(1)
	return &YouCodeDB{DataBase}
}

func (db *YouCodeDB) StartRefreshData() {

	// Refresh cache for channels
	channelsListsCached = db.getChannelsFromDB()
	ticker := time.NewTicker(1 * time.Hour)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				channelsListsCached = db.getChannelsFromDB()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
