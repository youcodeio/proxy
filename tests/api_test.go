package test

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type MyMockedObject struct {
	mock.Mock
}

func TestChannels(t *testing.T) {
	// create an instance of our test object
	testObj := new(MyMockedObject)
}
