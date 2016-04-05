package api

import (
	"encoding/json"
	"expvar"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"google.golang.org/api/youtube/v3"

	"github.com/gorilla/mux"
	"github.com/youcodeio/proxy/database"
	"github.com/youcodeio/proxy/utils"
)

// Two metrics, these are exposed by "magic" :)
// Number of calls to our server.
var lastSearch = expvar.NewString("youcodeio.last.search")
var numCalls = expvar.NewInt("youcodeio.counter.api.calls")

// NewRouter return a new mux Router
// https://groups.google.com/forum/#!msg/golang-nuts/Xs-Ho1feGyg/xg5amXHsM_oJ
func NewRouter(db *database.YouCodeDB) *mux.Router {
	r := mux.NewRouter()
	r.Handle("/channels", GetChannels(db))
	r.Handle("/search", GetQuery(db))
	r.Handle("/debug/vars", http.DefaultServeMux)
	return r
}

// GetChannels returns the list of channels used by YouCode
func GetChannels(db *database.YouCodeDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		numCalls.Add(1)

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

		numCalls.Add(1)

		query := r.URL.Query().Get("query")
		tuto := r.URL.Query().Get("istuts")
		if len(query) == 0 {
			http.Error(w, "Not enough args", http.StatusBadRequest)
			return
		}

		lastSearch.Set(query)

		channels := db.GetChannels()
		resultsChannel := make(chan []*youtube.SearchResult, database.MaxResults)

		var wg sync.WaitGroup
		var results []youtube.SearchResult
		for _, channel := range channels {
			if len(tuto) == 0 || strconv.FormatBool(channel.IsTuts) == tuto {
				log.Println("Querying", query, "on", channel.Name)
				wg.Add(1)
				go database.SearchOnChannel(query, channel.Ytid, resultsChannel, &wg)
			}
		}
		wg.Wait()
		log.Println("Fetching done")
		for index := 0; index < len(channels); index++ {
			for _, result := range <-resultsChannel {
				results = append(results, *result)
			}
		}

		// Sorting
		sort.Sort(utils.Channels(results))

		json, err := json.Marshal(results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	})
}
