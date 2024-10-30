package main

import (
	"log"
	"net/http"
	"os"
	"ws_server/router"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("ERROR LOADING .ENV")
	}

	port := os.Getenv("PORT")
	r := router.Router()

	http.ListenAndServe(port, r)
}
