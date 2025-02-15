package main

import (
	"doki.co.in/doki_real_time_service/hub"
	"doki.co.in/doki_real_time_service/payload"
	"fmt"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		// log.Printf("error loading env file: %v\n", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	userPoolID := os.Getenv("USER_POOL_ID")
	region := os.Getenv("REGION")
	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, userPoolID)

	jwks, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		// log.Fatalf("Failed to create JWK Set from resource at the given URL.\nError: %s", err)
	}

	fmt.Println("Doki real time service")
	// init payloads that can be received
	payload.InitPayload()
	newHub := hub.CreateHub(&jwks)
	http.HandleFunc("/ws", newHub.ServeWS)
	// log.Fatal(http.ListenAndServe(":"+port, nil))
}