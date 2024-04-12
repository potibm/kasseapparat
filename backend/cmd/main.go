package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := ":3000" // Default port number
	if len(os.Args) > 1 {
		port = ":" + os.Args[1] // Use the provided port number if available
	}

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	myHandler := func(w http.ResponseWriter, r *http.Request) {
		jsonData, err := os.ReadFile("./backend/data/products.json")
		if err != nil {
			log.Println("Error:", err)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(jsonData)
		}
	}

	http.HandleFunc("/api/products", myHandler)

	log.Println("Listening on " + port + "...")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
