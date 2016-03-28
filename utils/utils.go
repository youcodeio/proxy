package utils

import (
	"fmt"
	"time"

	"google.golang.org/api/youtube/v3"
)

type Channels []youtube.SearchResult

func (slice Channels) Len() int {
	return len(slice)
}

func (slice Channels) Less(i, j int) bool {
	layout := "2006-01-02T15:04:05.000Z"
	str1 := slice[i].Snippet.PublishedAt
	str2 := slice[i].Snippet.PublishedAt
	t1, err := time.Parse(layout, str1)
	t2, err := time.Parse(layout, str2)

	if err != nil {
		fmt.Println(err)
	}
	return t1.Before(t2)
}

func (slice Channels) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
