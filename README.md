# Sensor Monitoring Application

This application monitors sensor data and exposes metrics via Prometheus. It also checks the status of a microcontroller every 30 minutes.

## Features

- Monitor temperature, humidity, and feels-like data from sensors.
- Track door openings (stairs and garage).
- Expose metrics using Prometheus.
- Check microcontroller status periodically.

## Endpoints

| Endpoint         | Method | Description                                 |
|------------------|--------|---------------------------------------------|
| `/`              | POST   | Receive sensor data.                        |
| `/coloc-door`    | GET    | Open the stairs door.                       |
| `/coloc-door-garage` | GET | Open the garage door.                     |
| `/metrics`       | GET    | Retrieve Prometheus metrics.                |

## Metrics

| Metric                       | Description                                 |
|------------------------------|---------------------------------------------|
| `dht22_temperature_celsius`  | Temperature in Celsius.                     |
| `dht22_humidity_percent`     | Humidity in percentage.                     |
| `dht22_feelsLike_celsius`    | Feels-like temperature in Celsius.          |
| `coloc_stairs_counter`       | Number of times the stairs door has been opened. |
| `coloc_garage_counter`       | Number of times the garage door has been opened. |
| `coloc_esp_intercom_status`  | Status of the ESP intercom (1 if running, 0 if not). |

## Configuration

- **`COLOC_BASE_URL`**: Base URL for HTTP requests to the colocation system.

## Launch

```bash
go run main.go
```

The application will be accessible at http://localhost:8080.