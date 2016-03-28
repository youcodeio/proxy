// http://www.markjberger.com/testing-web-apps-in-golang/
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/youcodeio/proxy/database"
)

func TestStub(t *testing.T) {
	assert.True(t, true, "This is good. Canary test passing")
}

type MyMockedDB struct {
	mock.Mock
}

func (m *MyMockedDB) getChannelsFromDB() []database.Channel {
	args := m.Called()
	return args.Get(0).([]database.Channel)
}

/**
var (
	server *httptest.Server
)

func init() {
	mockedDB := new(MyMockedDB)
	server = httptest.NewServer(api.NewRouter(mockedDB)) //Creating new server
}



func TestSlash(t *testing.T) {

	assert := assert.New(t)
	request, err := http.NewRequest("GET", fmt.Sprintf("%s", server.URL), nil)

	res, err := http.DefaultClient.Do(request)

	assert.Nil(err)
	assert.Equal(res.StatusCode, 404, "Should be HTTP 404")
}

func TestChannelslist(t *testing.T) {

	assert := assert.New(t)
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/channels", server.URL), nil)

	res, err := http.DefaultClient.Do(request)

	assert.Nil(err)

	assert.Equal(res.StatusCode, 200, "Should be HTTP 200")

	var data []database.Channel

	defer res.Body.Close()
	// Read the content into a byte array
	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(err)

	err = json.Unmarshal(body, &data)
	assert.Nil(err)
}

type MyMockedDB struct {
	*database.YouCodeDB
	mock.Mock
}



func TestChannels(t *testing.T) {

	var channels []database.Channel
	ch := database.Channel{1, "name1", "id1", false}
	channels = append(channels, ch)
	ch = database.Channel{2, "name1", "id1", false}
	channels = append(channels, ch)

	mockedDB := new(MyMockedDB)
	mockedDB.On("getChannelsFromDB").Return(channels)

	r, err := http.NewRequest("GET", "employees/1", nil)
	w := httptest.NewRecorder()

	api.NewRouter(mockedDB).ServeHTTP(w, r)

}
**/
