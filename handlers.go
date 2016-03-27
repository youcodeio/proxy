package main

import (
	"encoding/json"
	"net/http"
)

// GetQuery return a result from a query
func GetQuery(w http.ResponseWriter, r *http.Request) {

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
