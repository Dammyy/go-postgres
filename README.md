# <h1>Time Sessions API/h1>

This is a simple API to store and retrieve time traccking sessions for https://github.com/Dammyy/tech-challenge-time  

## Getting Started

## Install
  go install

## Setting up the postgress DB
create a postgress dabase and provide the details in a .env file.
see .env.example for a list of values that have to be provided.

in your postgress DB, run this SQL command to create the table and columns 

```CREATE TABLE IF NOT EXISTS sessions(id SERIAL PRIMARY KEY, name text, time text, created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP)```

## start the app
 go run *.go