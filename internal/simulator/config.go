package simulator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Config holds the configuration settings for the simulation.
type Config struct {
	TotalReadings   int     `json:"totalReadings"`
	StartingTemp    float64 `json:"startingTemp"`
	MaxTempIncrease float64 `json:"maxTempIncrease"`
	TempFluctuation float64 `json:"tempFluctuation"`
	MinTemp         float64 `json:"minTemp"`
	MaxTemp         float64 `json:"maxTemp"`
	OutputFileName  string  `json:"outputFileName"`
	Simulate        bool    `json:"simulate"`
}

// Sensor represents the metadata for a sensor.
type Sensor struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	Version  string `json:"version"`
	Location string `json:"location"`
}

// SensorConfig holds the configuration and the sensor list.
type SensorConfig struct {
	Config  Config   `json:"config"`
	Sensors []Sensor `json:"sensors"`
}

// LoadConfigAndSensors loads the configuration and sensor metadata from a JSON file.
func LoadConfigAndSensors(filename string) (*SensorConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open configuration file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing configuration file: %v", err)
		}
	}()

	var sensorConfig SensorConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&sensorConfig); err != nil {
		return nil, fmt.Errorf("error decoding configuration JSON: %w", err)
	}

	if len(sensorConfig.Sensors) == 0 {
		return nil, fmt.Errorf("no sensors found in configuration")
	}

	return &sensorConfig, nil
}
