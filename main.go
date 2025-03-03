package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Data struct {
	Measurement string    `json:"measurement"`
	Host        string    `json:"host"`
	Value       string    `json:"value"`
	TypeValue   string    `json:"typeValue"`
	CreatedAt   time.Time `json:"createdAt"`
}

type metrics struct {
	temperature      prometheus.Gauge
	humidity         prometheus.Gauge
	feelsLike        prometheus.Gauge
	ColocDoorCounter prometheus.Counter
}

const COLOC_DOOR_BASE_URL = "http://10.0.0.2:3026"

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		temperature: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "dht22_temperature_celsius",
			Help: "Temperature from DHT22 sensor in Celsius.",
		}),
		humidity: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "dht22_humidity_percent",
			Help: "Humidity from DHT22 sensor in percentage.",
		}),
		feelsLike: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "dht22_feelsLike_celsius",
			Help: "Feels Like from DHT22 sensor in Celsius.",
		}),
		ColocDoorCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "coloc_door_counter",
			Help: "Number of times the door has been opened",
		}),
	}
	reg.MustRegister(m.temperature)
	reg.MustRegister(m.humidity)
	reg.MustRegister(m.feelsLike)
	reg.MustRegister(m.ColocDoorCounter)
	return m
}

func main() {

	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	// Create new metrics and register them using the custom registry.
	m := NewMetrics(reg)

	// HTTP Handler

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

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

		data.CreatedAt = time.Now()

		valueFormatted, err := strconv.ParseFloat(data.Value, 64)
		if err != nil {
			log.Fatal(err)
		}

		if data.Measurement == "temperature" {
			m.temperature.Set(valueFormatted)
		} else if data.Measurement == "humidity" {
			m.humidity.Set(valueFormatted)
		} else if data.Measurement == "realFeel" {
			m.feelsLike.Set(valueFormatted)
		}

		// Log the received data
		log.Printf("Received data: %+v\n", data)

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

	http.HandleFunc("/coloc-door", func(w http.ResponseWriter, r *http.Request) {
		// send a http request
		_, err := http.Get(COLOC_DOOR_BASE_URL + "/open")
		if err != nil {
			log.Fatal(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		_, res := w.Write([]byte("Data received successfully"))
		if res != nil {
			log.Printf("Error writing response: %v", res)
			return
		}

		m.ColocDoorCounter.Inc()
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
