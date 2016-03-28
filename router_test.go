package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	server *httptest.Server
)

func init() {
	server = httptest.NewServer(NewRouter()) //Creating new server
}

func TestChannelslist(t *testing.T) {

	assert := assert.New(t)
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/channels", server.URL), nil)

	res, err := http.DefaultClient.Do(request)

	assert.Nil(err)

	assert.Equal(res.StatusCode, 200, "Should be HTTP 200")

	var data []Channel

	defer res.Body.Close()
	// Read the content into a byte array
	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(err)

	err = json.Unmarshal(body, &data)
	assert.Nil(err)
}

func TestSlash(t *testing.T) {

	assert := assert.New(t)
	request, err := http.NewRequest("GET", fmt.Sprintf("%s", server.URL), nil)

	res, err := http.DefaultClient.Do(request)

	assert.Nil(err)
	assert.Equal(res.StatusCode, 404, "Should be HTTP 404")

}

func TestQuery(t *testing.T) {
	assert := assert.New(t)
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/query", server.URL), nil)

	res, err := http.DefaultClient.Do(request)
	assert.Nil(err)
	assert.Equal(res.StatusCode, 400, "Should not be able to query without query")

	request, err = http.NewRequest("GET", fmt.Sprintf("%s/query?query=\"golang\"", server.URL), nil)

	res, err = http.DefaultClient.Do(request)

	assert.Equal(res.StatusCode, 200, "Should be HTTP 200")

}
