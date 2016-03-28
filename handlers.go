package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"google.golang.org/api/youtube/v3"
)

// GetQuery return a result from a query
func GetQuery(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("query")
	if len(query) == 0 {
		http.Error(w, "Not enough args", http.StatusBadRequest)
		return
	}

	channels, err := getChannels()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resultsChannel := make(chan []youtube.SearchResultSnippet, MaxResults)

	var wg sync.WaitGroup
	var results []youtube.SearchResultSnippet
	for _, channel := range channels {
		wg.Add(1)
		go searchOnChannel(query, channel.Ytid, resultsChannel, &wg)
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
}

// Channelslist Return list of channels
func Channelslist(w http.ResponseWriter, r *http.Request) {

	channels, err := getChannels()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json, err := json.Marshal(channels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
