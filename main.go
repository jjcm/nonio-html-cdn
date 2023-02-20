package main

import (
	"fmt"
	"net/http"
	"os"
	"soci-html-cdn/config"
	"soci-html-cdn/route"
)

func setupRoutes(settings *config.Config) {
	http.Handle("/", http.FileServer(http.Dir("./files/html")))
	tmpFS := http.FileServer(http.Dir("./files/temp-html"))
	http.Handle("/temp/", http.StripPrefix("/temp/", tmpFS))
	http.HandleFunc("/upload", route.UploadFile)
	http.HandleFunc("/move", route.MoveFolder)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = settings.Port
		if port == "" {
			port = "4205"
		}
	}

	fmt.Printf("Listening on %v\n", port)
	http.ListenAndServe(":"+port, nil)
}

func main() {
	// parse the config file
	if err := config.ParseJSONFile("./config.json", &config.Settings); err != nil {
		panic(err)
	}
	// validate the config file
	if err := config.Settings.Validate(); err != nil {
		panic(err)
	}

	fmt.Println("Starting user-submitted html server...")
	setupRoutes(&config.Settings)
}
