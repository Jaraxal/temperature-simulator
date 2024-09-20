package test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"temperature-simulator/internal/simulator"
)

func TestLoadConfigAndSensors(t *testing.T) {
	// Valid configuration JSON content.
	validConfigJSON := `{
        "config": {
            "totalReadings": 5,
            "startingTemp": 20.0,
            "maxTempIncrease": 30.0,
            "tempFluctuation": 3.0,
            "minTemp": -50.0,
            "maxTemp": 100.0,
            "outputFileName": "temperature_readings",
            "outputFormat": "json",
            "simulate": true
        },
        "sensors": [
            {"name": "SensorA", "id": "001", "version": "v1.0", "location": "LocationA"},
            {"name": "SensorB", "id": "002", "version": "v1.1", "location": "LocationB"}
        ]
    }`

	// Create a temporary file with valid configuration JSON.
	validFile, err := os.CreateTemp("", "valid_config_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(validFile.Name())

	if _, err := validFile.Write([]byte(validConfigJSON)); err != nil {
		validFile.Close()
		t.Fatal(err)
	}
	validFile.Close()

	// Test loading valid configuration and sensors.
	sensorConfig, err := simulator.LoadConfigAndSensors(validFile.Name())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(sensorConfig.Sensors) != 2 {
		t.Errorf("Expected 2 sensors, got %d", len(sensorConfig.Sensors))
	}

	// Invalid configuration JSON content.
	invalidConfigJSON := `invalid json`

	// Create a temporary file with invalid configuration JSON.
	invalidFile, err := os.CreateTemp("", "invalid_config_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(invalidFile.Name())

	if _, err := invalidFile.Write([]byte(invalidConfigJSON)); err != nil {
		invalidFile.Close()
		t.Fatal(err)
	}
	invalidFile.Close()

	// Test loading invalid configuration.
	_, err = simulator.LoadConfigAndSensors(invalidFile.Name())
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

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
		t.Errorf("Expected no error, got %v", err)
	}

	expectedDataPoints := config.TotalReadings * len(sensors)
	if len(data) != expectedDataPoints {
		t.Errorf("Expected %d data points, got %d", expectedDataPoints, len(data))
	}

	for _, reading := range data {
		tempValue := float64(reading.Temperature)
		if tempValue < config.MinTemp || tempValue > config.MaxTemp {
			t.Errorf("Temperature out of bounds: %.2f", reading.Temperature)
		}
	}
}

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

	tmpfile, err := os.CreateTemp("", "test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	if err := simulator.SaveToJSON(data, tmpfile.Name()); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Read the file and check content.
	contentBytes, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	content := string(contentBytes)

	// Split the content into lines.
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
