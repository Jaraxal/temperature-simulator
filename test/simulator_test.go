// Package test contains unit tests for the simulator package, which includes
// functionality for loading sensor configurations, generating temperature readings,
// and saving data to JSON files.
package test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"temperature-simulator/internal/simulator"
)

// TestLoadConfigAndSensors tests the loading of sensor configurations from a JSON file.
// It verifies that the function correctly loads valid configurations and handles invalid
// file paths as expected.
func TestLoadConfigAndSensors(t *testing.T) {
	// Define the path to the external JSON configuration file.
	configFilePath := filepath.Join("..", "configs", "test_sensors.json")

	// Test loading valid configuration and sensors.
	sensorConfig, err := simulator.LoadConfigAndSensors(configFilePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(sensorConfig.Sensors) != 2 {
		t.Errorf("Expected 2 sensors, got %d", len(sensorConfig.Sensors))
	}

	// Test with an invalid configuration file path.
	invalidConfigFilePath := filepath.Join("..", "configs", "nonexistent.json")

	// Test loading invalid configuration.
	_, err = simulator.LoadConfigAndSensors(invalidConfigFilePath)
	if err == nil {
		t.Error("Expected error for invalid configuration file path, got nil")
	}
}

// TestGenerateTemperatureReadings tests the generation of temperature readings
// for a given set of sensors and configuration. It verifies that the correct number
// of readings are generated, and that the temperatures fall within the expected range.
func TestGenerateTemperatureReadings(t *testing.T) {
	sensors := []simulator.Sensor{
		{
			Name:     "SensorA",
			ID:       "001",
			Version:  "v1.0",
			Location: "LocationA",
		},
		{
			Name:     "SensorB",
			ID:       "002",
			Version:  "v1.1",
			Location: "LocationB",
		},
	}

	config := simulator.Config{
		TotalReadings:   10,
		StartingTemp:    20.0,
		MaxTempIncrease: 30.0,
		TempFluctuation: 3.0,
		MinTemp:         -10.0,
		MaxTemp:         50.0,
		Simulate:        true, // Use simulate mode to avoid actual sleeping
	}

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
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the number of generated readings matches expectations.
	expectedDataPoints := config.TotalReadings * len(sensors)
	if len(data) != expectedDataPoints {
		t.Errorf("Expected %d data points, got %d", expectedDataPoints, len(data))
	}

	// Verify the generated temperatures are within the allowed range.
	for _, reading := range data {
		tempValue := float64(reading.Temperature)
		if tempValue < config.MinTemp || tempValue > config.MaxTemp {
			t.Errorf("Temperature out of bounds: %.2f", reading.Temperature)
		}
	}
}

// TestSaveToJSON tests saving temperature readings to a JSON file.
// It verifies that the data is correctly written to the file in the expected format
// and that each line of the file corresponds to a valid JSON object.
func TestSaveToJSON(t *testing.T) {
	data := []simulator.TemperatureReading{
		{
			Time:        "2023-10-01 12:00:00",
			Temperature: simulator.Temperature(25.5),
			Sensor: simulator.Sensor{
				Name:     "SensorA",
				ID:       "001",
				Version:  "v1.0",
				Location: "LocationA",
			},
		},
		{
			Time:        "2023-10-01 12:01:00",
			Temperature: simulator.Temperature(26.0),
			Sensor: simulator.Sensor{
				Name:     "SensorB",
				ID:       "002",
				Version:  "v1.1",
				Location: "LocationB",
			},
		},
	}

	// Create a temporary file for testing.
	tmpfile, err := os.CreateTemp("", "test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Save the data to the JSON file.
	if err := simulator.SaveToJSON(data, tmpfile.Name()); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Read the file and check the content.
	contentBytes, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	content := string(contentBytes)

	// Split the content into lines and verify the number of lines matches the data.
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) != len(data) {
		t.Errorf("Expected %d lines, got %d", len(data), len(lines))
	}

	// Check each line is valid JSON and matches the expected data.
	for i, line := range lines {
		var reading simulator.TemperatureReading
		if err := json.Unmarshal([]byte(line), &reading); err != nil {
			t.Errorf("Error unmarshaling line %d: %v", i+1, err)
		}

		// Compare the readings, considering the temperature formatting.
		if reading.Time != data[i].Time {
			t.Errorf("Time mismatch on line %d.\nExpected: %s\nGot: %s", i+1, data[i].Time, reading.Time)
		}
		if reading.Sensor != data[i].Sensor {
			t.Errorf("Sensor mismatch on line %d.\nExpected: %+v\nGot: %+v", i+1, data[i].Sensor, reading.Sensor)
		}

		// Compare the temperatures, formatted to two decimal places.
		expectedTemp := fmt.Sprintf("%.2f", data[i].Temperature)
		actualTemp := fmt.Sprintf("%.2f", reading.Temperature)
		if expectedTemp != actualTemp {
			t.Errorf("Temperature mismatch on line %d.\nExpected: %s\nGot: %s", i+1, expectedTemp, actualTemp)
		}
	}
}
