package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"temperature-simulator/internal/simulator"
)

// SensorConfig holds the configuration for sensors.
type SensorConfig struct {
	LogLevel        string             `json:"log_level"`
	LogOutput       string             `json:"log_output"`
	TotalReadings   int                `json:"total_readings"`
	StartingTemp    float64            `json:"starting_temp"`
	MaxTempIncrease float64            `json:"max_temp_increase"`
	TempFluctuation float64            `json:"temp_fluctuation"`
	MinTemp         float64            `json:"min_temp"`
	MaxTemp         float64            `json:"max_temp"`
	Simulate        bool               `json:"simulate"`
	Sensors         []simulator.Sensor `json:"sensors"`
}

// GenerateTemperatureReadings generates temperature readings based on the provided configuration.
func GenerateTemperatureReadings(w http.ResponseWriter, r *http.Request) {
	var config SensorConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Setup logger based on the log level and output destination.
	if err := simulator.SetupLogger(config.LogLevel, config.LogOutput); err != nil {
		http.Error(w, "Error setting up logger: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Starting temperature simulator...")

	sensors := config.Sensors
	log.Printf("Loaded configuration: %+v", config)
	log.Printf("Loaded %d sensors", len(sensors))

	// Generate temperature readings.
	log.Println("Generating temperature readings...")
	data, err := simulator.GenerateTemperatureReadings(
		sensors,
		config.TotalReadings,
		config.StartingTemp,
		config.MaxTempIncrease,
		config.TempFluctuation,
		config.MinTemp,
		config.MaxTemp,
		config.Simulate,
	)
	if err != nil {
		http.Error(w, "Error generating temperature readings: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Generated %d temperature readings", len(data))

	// Return generated temperature readings as JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func main() {
	// Parse command-line flags for configuration file, log level, log output, and output file.
	sensorConfigFile := flag.String("sensor_config", "configs/sensors.json", "Path to the sensor configuration JSON file")
	logLevel := flag.String("log_level", "info", "Log level (debug, info, warn, error)")
	logOutput := flag.String("log_output", "", "Log output ('stdout' or file path), overrides config file log path")
	outputFile := flag.String("output", "", "Output file for temperature readings, overrides config file output file")
	flag.Parse()

	// Load the configuration and sensors from the JSON file.
	sensorConfig, err := simulator.LoadConfigAndSensors(*sensorConfigFile)
	if err != nil {
		log.Fatalf("Error loading configuration and sensors: %v", err)
	}

	// Use the log output from the config if the command-line flag is not provided.
	config := sensorConfig.Config
	if *logOutput == "" {
		*logOutput = config.LogFilePath
		if *logOutput == "" {
			*logOutput = "stdout" // Default to stdout if not specified in either place.
		}
	}

	http.HandleFunc("/generate-temperature-readings", GenerateTemperatureReadings)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
