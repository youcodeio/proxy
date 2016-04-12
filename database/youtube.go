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
	// DefaultMaxResults of Youtube API
	DefaultMaxResults = int64(25)
	service           *youtube.Service
)

func SearchOnChannel(q string, channel string, resultChannel chan []*youtube.SearchResult, wg *sync.WaitGroup, MaxResults int64) {

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

	if (len(response.Items)) != 0 {
		log.Println("Pushing", len(response.Items), "for", q, "from", channel)
		resultChannel <- response.Items
	}
	wg.Done()
}

//ChannelInfo returns info about a channel
func ChannelInfo(channelID string, resultChan chan *youtube.ChannelListResponse) {

	call := service.Channels.List("id,snippet,contentDetails,statistics").Id(channelID)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
	}
	resultChan <- response
}

func init() {
	var err error
	service, err = youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}
}
