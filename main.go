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
	a.Initialize(os.Getenv("DATABASE_URL"))

	port := os.Getenv("PORT")

    if port == "" {
        port = "8010"
	}
	
	a.Run(":" + port)
}
