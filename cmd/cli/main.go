package main

import (
	"flag"
	"log"

	"temperature-simulator/internal/simulator"
)

// main is the entry point of the temperature simulator application.
// It loads the sensor configuration, generates temperature readings,
// and saves the results in a JSON file.
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

	// Use the output file from the command-line flag, if provided, otherwise use the one from the config.
	if *outputFile != "" {
		config.OutputFileName = *outputFile
	}

	// Setup logger based on the log level and output destination.
	if err := simulator.SetupLogger(*logLevel, *logOutput); err != nil {
		log.Fatalf("Error setting up logger: %v", err)
	}

	log.Printf("Starting temperature simulator...")

	sensors := sensorConfig.Sensors
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
		log.Fatalf("Error generating temperature readings: %v", err)
	}
	log.Printf("Generated %d temperature readings", len(data))

	// Save generated temperature readings to the output file.
	log.Printf("Saving temperature readings to %s", config.OutputFileName)
	if err := simulator.SaveToJSON(data, config.OutputFileName); err != nil {
		log.Fatalf("Error saving to JSON: %v", err)
	}

	log.Println("Temperature simulation completed successfully.")
}
