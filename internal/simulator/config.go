package simulator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// SetupLogger configures the global logger based on the specified log level and output destination.
// The log level can be one of: "debug", "info", "warn", "error".
// The log output can be either "stdout" or a file path specified via command-line or configuration.
//
// Parameters:
//   - logLevel: The desired log level for the application.
//   - logOutput: The destination for the logs, either "stdout" or a file path.
//
// Returns:
//   - An error if the log level or log output is invalid, or nil if successful.
func SetupLogger(logLevel, logOutput string) error {
	// Determine the log output destination (stdout or a file).
	var output *os.File
	if logOutput == "stdout" {
		output = os.Stdout
	} else {
		// Open or create the log file.
		var err error
		output, err = os.OpenFile(logOutput, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("unable to open log file: %w", err)
		}
	}

	log.SetOutput(output)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // Add timestamp and file info to logs.

	// Set the log level.
	switch strings.ToLower(logLevel) {
	case "debug":
		log.SetPrefix("DEBUG: ")
	case "info":
		log.SetPrefix("INFO: ")
	case "warn":
		log.SetPrefix("WARN: ")
	case "error":
		log.SetPrefix("ERROR: ")
	default:
		return fmt.Errorf("unknown log level: %s", logLevel)
	}

	log.Printf("Logger initialized with level: %s, output: %s", logLevel, logOutput)
	return nil
}

// Config holds the configuration settings for the temperature simulation, including log file path.
// This struct defines the core parameters for running the simulation, such as the number of readings,
// initial temperature, temperature fluctuations, and the simulation mode.
type Config struct {
	TotalReadings   int     `json:"totalReadings"`   // Number of temperature readings to generate.
	StartingTemp    float64 `json:"startingTemp"`    // Initial temperature for all sensors at the start of the simulation.
	MaxTempIncrease float64 `json:"maxTempIncrease"` // Maximum temperature increase allowed during the increase period.
	TempFluctuation float64 `json:"tempFluctuation"` // The maximum random fluctuation to be applied to the temperature.
	MinTemp         float64 `json:"minTemp"`         // The minimum allowable temperature value.
	MaxTemp         float64 `json:"maxTemp"`         // The maximum allowable temperature value.
	OutputFileName  string  `json:"outputFileName"`  // Name of the file where simulation results will be saved.
	Simulate        bool    `json:"simulate"`        // If true, the simulation runs over real time; otherwise, it runs as fast as possible.
	LogFilePath     string  `json:"logFilePath"`     // Path to the log file, if not provided via command-line.
}

// Sensor holds metadata information about a specific sensor used in the simulation.
// Each sensor is identified by its name, ID, version, and physical location.
type Sensor struct {
	Name     string `json:"name"`     // Human-readable name of the sensor (e.g., "Sensor A").
	ID       string `json:"id"`       // Unique identifier for the sensor.
	Version  string `json:"version"`  // Version information about the sensor.
	Location string `json:"location"` // Physical location or placement of the sensor.
}

// SensorConfig represents the complete configuration for the simulation.
// It includes the global simulation configuration and a list of sensors that will
// generate temperature readings.
type SensorConfig struct {
	Config  Config   `json:"config"`  // Global simulation configuration settings.
	Sensors []Sensor `json:"sensors"` // List of sensors to simulate.
}

// LoadConfigAndSensors loads the simulation configuration and sensor metadata from a JSON file.
// It reads the configuration file, decodes it into a SensorConfig struct, and returns the struct.
//
// Parameters:
//   - filename: The path to the configuration file containing the simulation and sensor settings.
//
// Returns:
//   - A pointer to a SensorConfig struct populated with the configuration and sensors.
//   - An error if the file cannot be opened or if the JSON is invalid.
//
// This function uses buffered reading for efficiency, especially with larger configuration files.
// It will return an error if no sensors are found in the configuration or if the JSON format is incorrect.
func LoadConfigAndSensors(filename string) (*SensorConfig, error) {
	// Open the configuration file for reading.
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening configuration file: %v", err)
		return nil, fmt.Errorf("unable to open configuration file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing configuration file: %v", err)
		}
	}()

	// Use a buffered reader for efficient reading of the file contents.
	reader := bufio.NewReader(file)

	// Decode the JSON configuration into a SensorConfig struct.
	var sensorConfig SensorConfig
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&sensorConfig); err != nil {
		log.Printf("Error decoding JSON configuration: %v", err)
		return nil, fmt.Errorf("error decoding configuration JSON: %w", err)
	}

	// Ensure that at least one sensor is defined in the configuration.
	if len(sensorConfig.Sensors) == 0 {
		log.Printf("No sensors found in configuration")
		return nil, fmt.Errorf("no sensors found in configuration")
	}

	// Log a message after loading the sensors successfully
	log.Printf("Loaded %d sensors from configuration", len(sensorConfig.Sensors))

	return &sensorConfig, nil
}
