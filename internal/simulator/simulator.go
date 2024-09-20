package simulator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

// Temperature is a custom type for temperature values.
type Temperature float64

// MarshalJSON implements the json.Marshaler interface for Temperature.
func (t Temperature) MarshalJSON() ([]byte, error) {
	formattedTemp := fmt.Sprintf("%.2f", t)
	return []byte(formattedTemp), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for Temperature.
func (t *Temperature) UnmarshalJSON(b []byte) error {
	var temp float64
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}
	*t = Temperature(temp)
	return nil
}

// TemperatureReading represents a single temperature reading.
type TemperatureReading struct {
	Time        string      `json:"time"`
	Temperature Temperature `json:"temperature"`
	Sensor      Sensor      `json:"sensor"`
}

const (
	timeFormat            = "2006-01-02 15:04:05"
	readingsPerHour       = 60
	increasePeriodMinutes = 5
)

// GenerateTemperatureReadings generates temperature readings for the given sensors.
func GenerateTemperatureReadings(
	sensors []Sensor,
	totalReadings int,
	startingTemp, maxTempIncrease, tempFluctuation, minTemp, maxTemp float64,
	simulate bool,
) ([]TemperatureReading, error) {
	// Initialize temperature for each sensor.
	sensorTemps := make(map[string]float64)
	for _, sensor := range sensors {
		sensorTemps[sensor.Name] = startingTemp
	}

	var data []TemperatureReading

	increaseAmountPerMinute := maxTempIncrease / float64(increasePeriodMinutes)

	// Initialize currentTime.
	var currentTime time.Time
	if simulate {
		currentTime = time.Now().UTC()
	}

	for loopCount := 0; loopCount < totalReadings; loopCount++ {
		// Use the configurable timeSleep function.
		if !simulate {
			time.Sleep(60 * time.Second)
		}

		// Update currentTime based on simulation mode.
		if simulate {
			currentTime = currentTime.Add(60 * time.Second)
		} else {
			currentTime = time.Now().UTC()
		}

		// Determine if increase phase is active.
		increasePhase := false
		if loopCount%readingsPerHour < increasePeriodMinutes {
			increasePhase = true
		}

		for _, sensor := range sensors {
			temp := sensorTemps[sensor.Name]

			// Apply normal temperature fluctuation.
			fluctuation := rand.Float64()*2*tempFluctuation - tempFluctuation
			temp += fluctuation

			// Apply temperature increase during increase phase.
			if increasePhase {
				temp += increaseAmountPerMinute
			}

			// Ensure temperature is within bounds.
			if temp < minTemp {
				temp = minTemp
			} else if temp > maxTemp {
				temp = maxTemp
			}

			// Update the sensor's temperature.
			sensorTemps[sensor.Name] = temp

			// Create and store the reading.
			reading := TemperatureReading{
				Time:        currentTime.Format(timeFormat),
				Temperature: Temperature(temp),
				Sensor:      sensor,
			}
			data = append(data, reading)
		}
	}

	return data, nil
}

// SaveToJSON saves the temperature readings to a NDJSON file.
func SaveToJSON(data []TemperatureReading, filename string) error {
	// Create file.
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating JSON file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing JSON file: %v", err)
		}
	}()

	// Create a buffered writer for efficiency.
	writer := bufio.NewWriter(file)
	defer func() {
		if err := writer.Flush(); err != nil {
			log.Printf("Error flushing JSON writer: %v", err)
		}
	}()

	// Encode each reading as a single line JSON object.
	for _, reading := range data {
		jsonData, err := json.Marshal(reading)
		if err != nil {
			return fmt.Errorf("error encoding JSON data: %w", err)
		}
		// Write the JSON object followed by a newline.
		if _, err := writer.Write(jsonData); err != nil {
			return fmt.Errorf("error writing JSON data: %w", err)
		}
		if err := writer.WriteByte('\n'); err != nil {
			return fmt.Errorf("error writing newline: %w", err)
		}
	}

	fmt.Printf("Data successfully saved to %s\n", filename)
	return nil
}
