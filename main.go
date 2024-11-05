package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"ws_server/router"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("ERROR LOADING .ENV")
		fmt.Printf("ERROR")
	}

	port := os.Getenv("PORT")
	redisAddr := os.Getenv("REDIS_ADDR")
	r := router.Router(redisAddr)

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal("ERROR WHILE SERVING HTTP")
	}
}
