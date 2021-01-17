# <h1>Time Sessions API</h1>
  This is a simple API to store and retrieve time traccking sessions for https://github.com/Dammyy/tech-challenge-time  

## Getting Started
  Deployed api - https://time-sessions-go-api.herokuapp.com/sessions  
  To setup locally, follow the instructions below

## Setting up the postgres DB
  Create a postgres database and provide the details in a `.env` file.  
  See `.env.example` for a list of values that have to be provided.

## Install
  Make sure you have Go version 1.12 or newer installed.  
  Run `go build -o bin/time-sessions -v`

## Start the app
 `bin/time-sessions`
 
## API 
  - GET all sessions `/sessions`
  - POST session `/session`
  - GET session `/sessions/:id`
  - PUT session `/sessions/:id`
  - DELETE session `/sessions/:id`
