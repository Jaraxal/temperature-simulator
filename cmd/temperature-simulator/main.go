package main

import (
	"flag"
	"log"

	"temperature-simulator/internal/simulator"
)

func init() {
	// Set log output to stderr and disable timestamps.
	log.SetFlags(0)
}

func main() {
	// Path to the sensor configuration JSON file.
	sensorConfigFile := flag.String("sensor_config", "configs/sensors.json", "Path to the sensor configuration JSON file")
	flag.Parse()

	// Load configuration and sensors from JSON file.
	sensorConfig, err := simulator.LoadConfigAndSensors(*sensorConfigFile)
	if err != nil {
		log.Fatalf("Error loading configuration and sensors: %v", err)
	}

	config := sensorConfig.Config
	sensors := sensorConfig.Sensors

	// Generate temperature readings.
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
		log.Fatalf("Error generating temperature readings: %v", err)
	}

	// Save data to file.
	outputFileName := config.OutputFileName
	if outputFileName == "" {
		outputFileName = "output/temperature-readings"
	}

	if err := simulator.SaveToJSON(data, outputFileName+".json"); err != nil {
		log.Fatalf("Error saving to JSON: %v", err)
	}
}
