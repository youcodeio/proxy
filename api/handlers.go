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
	r := mux.NewRouter().StrictSlash(true)
	r.Handle("/channels", GetChannels(db))
	r.Handle("/channels/{channel}/lastVideos", GetChannelLatestVideos())
	r.Handle("/channels/{channel}/info", GetChannelInfo())
	r.Handle("/search", GetQuery(db))
	r.Handle("/debug/vars", http.DefaultServeMux)
	return r
}

// GetChannels returns the list of channels used by YouCode
func GetChannels(db *database.YouCodeDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}

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

//GetChannelLatestVideos returns the latest videos from a channel
func GetChannelLatestVideos() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}
		numCalls.Add(1)

		var videoSize int64
		var err error
		size := r.URL.Query().Get("size")
		if len(size) == 0 {
			videoSize = 2
		} else {
			videoSize, err = strconv.ParseInt(size, 10, 64)
			if err != nil {
				http.Error(w, "size is not a number", http.StatusBadRequest)
				return
			}
		}

		vars := mux.Vars(r)
		id := vars["channel"]

		if len(id) == 0 {
			http.Error(w, "Not enough args", http.StatusBadRequest)
			return
		}

		resultsChannel := make(chan []*youtube.SearchResult, videoSize)

		var wg sync.WaitGroup
		var results []youtube.SearchResult

		wg.Add(1)
		go database.SearchOnChannel("", id, resultsChannel, &wg, videoSize)

		wg.Wait()
		close(resultsChannel)

		log.Println("Size of resultsChannels", len(resultsChannel))

		for _, result := range <-resultsChannel {
			// log.Println("Pushing ", *result.Id)
			results = append(results, *result)
		}

		log.Println("Sorting results...")
		// Sorting
		sort.Sort(utils.Channels(results))
		log.Println("Done")

		json, err := json.Marshal(results)
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
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}
		numCalls.Add(1)

		query := r.URL.Query().Get("query")
		tuto := r.URL.Query().Get("istuts")
		if len(query) == 0 {
			http.Error(w, "Not enough args", http.StatusBadRequest)
			return
		}

		lastSearch.Set(query)

		channels := db.GetChannels()
		resultsChannel := make(chan []*youtube.SearchResult, database.DefaultMaxResults*int64(len(channels)))

		var wg sync.WaitGroup
		var results []youtube.SearchResult
		for _, channel := range channels {
			if len(tuto) == 0 || strconv.FormatBool(channel.IsTuts) == tuto {
				log.Println("Querying", query, "on", channel.Name)
				wg.Add(1)
				go database.SearchOnChannel(query, channel.Ytid, resultsChannel, &wg, database.DefaultMaxResults)
			}
		}

		wg.Wait()
		close(resultsChannel)

		log.Println("Size of resultsChannels", len(resultsChannel))

		for index := 0; index < len(channels); index++ {
			// Because we closed the channel above,
			// the iteration terminates after receiving the events
			for _, result := range <-resultsChannel {
				// log.Println("Pushing ", *result.Id)
				results = append(results, *result)
			}
		}

		log.Println("Sorting results...")
		// Sorting
		sort.Sort(utils.Channels(results))
		log.Println("Done")

		json, err := json.Marshal(results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	})
}

func GetChannelInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}

		vars := mux.Vars(r)
		id := vars["channel"]
		resultChannel := make(chan *youtube.ChannelListResponse)

		log.Println("Querying channel", id, "for info")

		go database.ChannelInfo(id, resultChannel)

		json, err := json.Marshal(<-resultChannel)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	})
}
