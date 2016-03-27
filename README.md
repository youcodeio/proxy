# proxy [![Build Status](https://travis-ci.org/youcodeio/proxy.svg?branch=master)](https://travis-ci.org/youcodeio/proxy) [![Go Report Card](https://goreportcard.com/badge/github.com/youcodeio/proxy)](https://goreportcard.com/report/github.com/youcodeio/proxy) [![codebeat badge](https://codebeat.co/badges/40d6e665-663d-43db-8380-b58755d8a4aa)](https://codebeat.co/projects/github-com-youcodeio-proxy) [![GoDoc](https://godoc.org/github.com/youcodeio/proxy?status.svg)](https://godoc.org/github.com/youcodeio/proxy)
Repo for youcode.io V2.0. Endpoint for the API

# DB model
        CREATE TABLE CHANNELS(
                ID             SERIAL PRIMARY KEY,
                NAME           TEXT      NOT NULL,
                YTID           TEXT      NOT NULL,
                ISTUTS         BOOLEAN   NOT NULL
    );

# Insert data into DB
    INSERT INTO CHANNELS (NAME, YTID, ISTUTS) VALUES ('Google Developers', 'UC_x5XG1OV2P6uZZ5FSM9Ttw', false);
