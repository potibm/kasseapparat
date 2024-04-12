package main

import (
	"encoding/json"
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
		data := map[string]interface{}{}
		data["message"] = "Hello, World!"
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Println("Error:", err)
		} else {
			w.Header().Set("Content-Type", "application/json")
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
