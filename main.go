package main

import (
	"doki.co.in/doki_real_time_service/hub"
	"fmt"
	"log"
	"net/http"
)

func main() {

	fmt.Println("Doki real time service")

	newHub := hub.CreateHub()
	http.HandleFunc("/ws", newHub.ServeWS)
	log.Fatal(http.ListenAndServe(":8080", nil))
}