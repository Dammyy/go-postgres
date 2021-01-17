// main.go

package main

import (
	"os"
	"log"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
        log.Fatalf("Error loading .env file")
	}
	
	a := App{}
	a.Initialize(
		os.Getenv("APP_DATABASE_USERNAME"),
		os.Getenv("APP_DATABASE_PASSWORD"),
		os.Getenv("APP_DATABASE_NAME"))

    port := os.Getenv("PORT")
    if err != nil {
        port = "8010"
    } 

	a.Run(":" + port)
}
