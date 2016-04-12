# proxy [![Build Status](https://travis-ci.org/youcodeio/proxy.svg?branch=master)](https://travis-ci.org/youcodeio/proxy) [![Go Report Card](https://goreportcard.com/badge/github.com/youcodeio/proxy)](https://goreportcard.com/report/github.com/youcodeio/proxy) [![codebeat badge](https://codebeat.co/badges/40d6e665-663d-43db-8380-b58755d8a4aa)](https://codebeat.co/projects/github-com-youcodeio-proxy) [![GoDoc](https://godoc.org/github.com/youcodeio/proxy?status.svg)](https://godoc.org/github.com/youcodeio/proxy)
Repo for youcode.io V2.0. Endpoint for the API

# API Doc

endpoint: https://youcode-backend.cleverapps.io/

| Verb | URL | Description |
| ------------ | ------------ | ------------ |
| GET | /channels | Return all the channels handled by YouCode |
| GET | /channels/{channel_ID}/lastVideos | Returns the last 2 videos from a channel. Optional query: ?size to select video number |
| GET | /channels/{channel_ID}/info| Returns information about the channel |
| GET | /search?query=Go&istuts=true| Returns all the videos matching Golang from the selected channels. istuts is optionnal(default to false) |
| GET | /debug/vars | Expvar mode. Exposing go's metrics and also last research and number of calls for the API
