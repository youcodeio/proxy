package database

import (
	"log"
	"net/http"
	"os"
	"sync"

	"google.golang.org/api/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
)

var (
	client = &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}
	developerKey = os.Getenv("KEY_YT")
	// MaxResults of Youtube API
	MaxResults = 25
)

func SearchOnChannel(q string, channel string, resultChannel chan []youtube.SearchResultSnippet, wg *sync.WaitGroup) {
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}
	// Make the API call to YouTube.
	call := service.Search.List("id,snippet").
		ChannelId(channel).
		Order("date").
		Type("video").
		Q(q).
		MaxResults(int64(MaxResults))
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
	}

	var results []youtube.SearchResultSnippet

	// Parsing result into a clean slice
	for _, item := range response.Items {
		results = append(results, *item.Snippet)
	}
	wg.Done()
	resultChannel <- results
}
