package main

import (
	"log"
	"net/http"

	"github.com/youcodeio/proxy/api"
	"github.com/youcodeio/proxy/database"
)

func main() {
	db := database.InitDB()

	// We are cahcing results from DB. We need to start a goroutine
	db.StartRefreshData()

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", api.NewRouter(db)))
}
