package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"google.golang.org/api/youtube/v3"

	"github.com/gorilla/mux"
	"github.com/youcodeio/proxy/database"
)

// NewRouter return a new mux Router
// https://groups.google.com/forum/#!msg/golang-nuts/Xs-Ho1feGyg/xg5amXHsM_oJ
func NewRouter(db *database.YouCodeDB) *mux.Router {
	r := mux.NewRouter()
	r.Handle("/channels", GetChannels(db))
	r.Handle("/query", GetQuery(db))
	return r
}

// GetChannels returns the list of channels used by YouCode
func GetChannels(db *database.YouCodeDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		channels := db.GetChannels()
		json, err := json.Marshal(channels)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	})
}

// GetQuery return the result of a query on all the channels available
func GetQuery(db *database.YouCodeDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		if len(query) == 0 {
			http.Error(w, "Not enough args", http.StatusBadRequest)
			return
		}

		channels := db.GetChannels()
		resultsChannel := make(chan []youtube.SearchResultSnippet, database.MaxResults)

		var wg sync.WaitGroup
		var results []youtube.SearchResultSnippet
		for _, channel := range channels {
			wg.Add(1)
			go database.SearchOnChannel(query, channel.Ytid, resultsChannel, &wg)
		}
		wg.Wait()
		log.Println("Fetching done")
		for index := 0; index < len(channels); index++ {
			for _, result := range <-resultsChannel {
				results = append(results, result)
			}
		}

		json, err := json.Marshal(results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	})
}
