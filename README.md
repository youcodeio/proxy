# proxy
Repo for youcode.io V2.0. Endpoint for the API

# db model
CREATE TABLE CHANNELS(
   ID             SERIAL PRIMARY KEY,
   NAME           TEXT      NOT NULL,
   YTID           TEXT      NOT NULL,
   ISTUTS         BOOLEAN   NOT NULL
);

# Insert
INSERT INTO CHANNELS (NAME, YTID, ISTUTS) VALUES ('Google Developers', 'UC_x5XG1OV2P6uZZ5FSM9Ttw', false);
