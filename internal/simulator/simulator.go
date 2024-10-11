// Package simulator provides tools for simulating temperature readings from various sensors
// and saving them in a structured format.
package simulator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Temperature is a custom type representing temperature values in the simulation.
// It allows for custom JSON marshaling and unmarshaling to handle temperature formatting.
type Temperature float64

// MarshalJSON formats Temperature values with two decimal places when encoding to JSON.
func (t Temperature) MarshalJSON() ([]byte, error) {
	formattedTemp := strconv.FormatFloat(float64(t), 'f', 2, 64)
	return []byte(formattedTemp), nil
}

// UnmarshalJSON parses JSON data to populate a Temperature value.
// It expects the JSON data to be a float64 and converts it to the Temperature type.
func (t *Temperature) UnmarshalJSON(b []byte) error {
	temp, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		return err
	}
	*t = Temperature(temp)
	return nil
}

// TemperatureReading represents a single temperature reading from a sensor.
// It contains the time of the reading, the temperature value, and sensor metadata.
type TemperatureReading struct {
	Time        string      `json:"time"`        // Time of the reading in UTC format.
	Temperature Temperature `json:"temperature"` // The measured temperature value.
	Sensor      Sensor      `json:"sensor"`      // Metadata about the sensor making the reading.
}

const (
	// timeFormat specifies the layout used for formatting timestamps in the simulation.
	timeFormat = "2006-01-02 15:04:05"

	// readingsPerHour defines how many readings are taken per hour.
	readingsPerHour = 60

	// increasePeriodMinutes defines how many minutes each hour the temperature is increased.
	increasePeriodMinutes = 5
)

// GenerateTemperatureReadings simulates temperature readings for the specified sensors.
// It generates `totalReadings` temperature readings for each sensor, starting from `startingTemp`.
// The function also models temperature fluctuation and optional temperature increases during a
// predefined period (`increasePeriodMinutes`).
//
// Parameters:
//   - sensors: List of Sensor objects for which readings are generated.
//   - totalReadings: Total number of readings to generate for each sensor.
//   - startingTemp: The initial temperature value for all sensors.
//   - maxTempIncrease: The maximum amount by which the temperature can increase during an increase phase.
//   - tempFluctuation: The maximum amount of random fluctuation applied to the temperature in each reading.
//   - minTemp: The minimum allowable temperature value.
//   - maxTemp: The maximum allowable temperature value.
//   - simulate: If true, simulates readings over time; otherwise, fast-forwards the simulation.
//
// Returns a slice of `TemperatureReading` objects and an error (if applicable).
func GenerateTemperatureReadings(
	sensors []Sensor,
	totalReadings int,
	startingTemp, maxTempIncrease, tempFluctuation, minTemp, maxTemp float64,
	simulate bool,
) ([]TemperatureReading, error) {

	// Log the start of temperature generation
	log.Printf("Starting temperature generation for %d sensors with %d readings each", len(sensors), totalReadings)

	// Initialize temperature values for each sensor.
	sensorTemps := make([]float64, len(sensors))
	for i := range sensors {
		sensorTemps[i] = startingTemp
	}

	// Preallocate data slice to avoid resizing in the loop.
	data := make([]TemperatureReading, 0, totalReadings*len(sensors))

	// Calculate the temperature increase per minute.
	increaseAmountPerMinute := maxTempIncrease / float64(increasePeriodMinutes)
	var currentTime time.Time
	if simulate {
		currentTime = time.Now().UTC()
	}

	// Create a random number generator with a seed based on the current time.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate temperature readings for the required number of readings.
	for loopCount := 0; loopCount < totalReadings; loopCount++ {
		if !simulate {
			// Sleep for 60 seconds between readings if real-time simulation is disabled.
			time.Sleep(60 * time.Second)
		}
		// Update the current time, depending on whether simulation is active.
		if simulate {
			currentTime = currentTime.Add(60 * time.Second)
		} else {
			currentTime = time.Now().UTC()
		}

		// Determine if we're in the temperature increase phase.
		increasePhase := loopCount%readingsPerHour < increasePeriodMinutes

		for i, sensor := range sensors {
			temp := sensorTemps[i]

			// Apply random temperature fluctuation.
			fluctuation := r.Float64()*2*tempFluctuation - tempFluctuation
			temp += fluctuation

			// Apply a temperature increase if in the increase phase.
			if increasePhase {
				temp += increaseAmountPerMinute
			}

			// Ensure the temperature is within the specified min/max range.
			if temp < minTemp {
				temp = minTemp
			} else if temp > maxTemp {
				temp = maxTemp
			}

			// Store the updated temperature back to the sensor.
			sensorTemps[i] = temp

			// Create a new reading with the updated temperature and current time.
			reading := TemperatureReading{
				Time:        currentTime.Format(timeFormat),
				Temperature: Temperature(temp),
				Sensor:      sensor,
			}
			data = append(data, reading)
		}
	}

	log.Printf("Completed temperature generation. Total readings generated: %d", len(data))
	return data, nil
}

// SaveToJSON writes the temperature readings to a file in NDJSON (newline-delimited JSON) format.
// Each line in the output file represents a single JSON object containing a temperature reading.
//
// Parameters:
//   - data: The temperature readings to write.
//   - filename: The name of the file to save the readings to.
//
// Returns an error if the file cannot be created or written to.
func SaveToJSON(data []TemperatureReading, filename string) error {
	// Create the output file for writing.
	log.Printf("Saving data to JSON file: %s", filename)
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return fmt.Errorf("error creating JSON file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Printf("Error closing JSON file: %v", cerr)
		}
	}()

	// Use a buffered writer for improved performance.
	writer := bufio.NewWriterSize(file, 4096)
	defer func() {
		if err := writer.Flush(); err != nil {
			log.Printf("Error flushing JSON writer: %v", err)
		}
	}()

	// Write each temperature reading to the file as a JSON object.
	for _, reading := range data {
		// Marshal the reading to JSON format.
		jsonData, err := json.Marshal(reading)
		if err != nil {
			log.Printf("Error encoding JSON: %v", err)
			return fmt.Errorf("error encoding JSON data: %w", err)
		}

		// Write the JSON data followed by a newline.
		if _, err := writer.Write(jsonData); err != nil {
			log.Printf("Error writing JSON data: %v", err)
			return fmt.Errorf("error writing JSON data: %w", err)
		}
		if err := writer.WriteByte('\n'); err != nil {
			log.Printf("Error writing newline: %v", err)
			return fmt.Errorf("error writing newline: %w", err)
		}
	}

	log.Printf("Data successfully saved to %s", filename)
	return nil
}
