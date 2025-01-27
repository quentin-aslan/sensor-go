package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Data struct {
	Measurement string `json:"measurement"`
	Host        string `json:"host"`
	Value       string `json:"value"`
	TypeValue   string `json:"typeValue"`
}

var allData []Data

func main() {

	http.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(allData); err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			log.Printf("Error encoding JSON: %v", err)
			return
		}
	}))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
			return
		}

		var data Data
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			log.Printf("Error decoding JSON: %v", err)
			return
		}
		defer r.Body.Close()

		// Log the received data
		log.Printf("Received data: %+v\n", data)
		allData = append(allData, data)

		// Réponse au client
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		_, res := w.Write([]byte("Data received successfully"))
		if res != nil {
			log.Printf("Error writing response: %v", res)
			return
		}
		fmt.Println(w, "Data received successfully")
	})

	// Lancer le serveur HTTP
	port := ":8080"
	server := &http.Server{
		Addr:           port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Server started on port %s\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
