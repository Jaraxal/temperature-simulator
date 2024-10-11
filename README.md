# Temperature Simulator

A Go-based temperature simulator that generates temperature readings for a set of sensors. The output of the simulator is intended to represent temperature data as one might expect to find in data center racks or build HVAC systems.  The ultimate goal is to load the data into a tool like Elasticsearch for use with Dashboards, Detection Rules and Machine Learning.

This project was created as part of an effort to learn Go, as a long time-time Perl, Python and Node.js programmer. I would appreciate constructive feedback.

## Table of Contents

- [Temperature Simulator](#temperature-simulator)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Installation](#installation)
    - [Prerequisites](#prerequisites)
    - [Clone the Repository](#clone-the-repository)
    - [Build the Project](#build-the-project)
  - [Usage](#usage)
    - [Running the Simulator](#running-the-simulator)
    - [Command-Line Options](#command-line-options)
  - [Configuration](#configuration)
    - [Example Configuration](#example-configuration)
    - [Configuration Parameters](#configuration-parameters)
    - [Sensors Configuration](#sensors-configuration)
  - [Directory Structure](#directory-structure)
  - [Testing](#testing)
  - [Useful Commands](#useful-commands)
    - [Build the project](#build-the-project-1)
    - [Run the project](#run-the-project)
    - [Clean the build and temp files](#clean-the-build-and-temp-files)
    - [Run tests](#run-tests)
    - [Format Go files](#format-go-files)
    - [Run linting](#run-linting)
    - [Install dependencies](#install-dependencies)
    - [Run everything (lint, fmt, test, build)](#run-everything-lint-fmt-test-build)
    - [Installing golangci-lint](#installing-golangci-lint)
  - [Contributing](#contributing)

## Features

- Simulate temperature readings for multiple sensors
- Configurable parameters for flexible simulations
- Easy-to-use command-line interface
- Unit tests included for reliability

## Installation

### Prerequisites

Go 1.20 or higher installed on your system. You can download it from the official website or use `brew` on Mac OS.

### Clone the Repository

```bash
git clone https://github.com/jaraxal/temperature-simulator.git
cd temperature-simulator
```

### Build the Project

```bash
go build -o ./bin/temperature-simulator ./cmd/temperature-simulator
```

## Usage

### Running the Simulator

By default, the simulator uses the configuration file located at configs/sensors.json.

```bash
./temperature-simulator
```

You can specify a different configuration file using the -sensor_config flag:

```bash
./temperature-simulator -sensor_config=path/to/your_config.json
```

### Command-Line Options

- `-sensor_config`: Path to the sensor configuration JSON file. Default is configs/sensors.json.
- `-log_output`: Specify where to write log output (stdout for terminal or a file path).
- `-log_level`: Log level (e.g., debug, info, warn, error).
- `-output`: Override the output file name specified in the configuration file.

## Configuration

The simulator is configured via a JSON file that specifies both the simulation parameters and the sensor metadata.

### Example Configuration

`configs/sensors.json`

```json
{
  "config": {
    "totalReadings": 5,
    "startingTemp": 20.0,
    "maxTempIncrease": 30.0,
    "tempFluctuation": 3.0,
    "minTemp": -50.0,
    "maxTemp": 100.0,
    "outputFileName": "output/temperature-readings.json",
    "simulate": true,
    "logFilePath": "logs/temperature-simulator.log"
  },
  "sensors": [
    {
      "name": "SensorA",
      "id": "001",
      "version": "v1.0",
      "location": "LocationA"
    },
    {
      "name": "SensorB",
      "id": "002",
      "version": "v1.1",
      "location": "LocationB"
    },
    {
      "name": "SensorC",
      "id": "003",
      "version": "v2.0",
      "location": "LocationC"
    }
  ]
}
```

### Configuration Parameters

- `totalReadings`: Number of readings to generate.
- `startingTemp`: The starting temperature for all sensors.
- `maxTempIncrease`: The maximum temperature increase during the increase phase.
- `tempFluctuation`: The maximum fluctuation in temperature per reading.
- `minTemp`: The minimum allowable temperature.
- `maxTemp`: The maximum allowable temperature.
- `outputFileName`: The file name of the json output file.
- `simulate`: If true, the simulator runs without actual time delays.
- `logFilePath`: The file name of the log file output.

### Sensors Configuration

Each sensor in the sensors array has the following fields:

- `name`: The name of the sensor.
- `id`: The unique identifier of the sensor.
- `version`: The version of the sensor hardware or firmware.
- `location`: The physical location of the sensor.

## Directory Structure

```go
temperature-simulator/
├── cmd/
│   └── temperature-simulator/
│       └── main.go
├── configs/
│   └── sensors.json
│   └── test_sensors.json
├── internal/
│   └── simulator/
│       ├── config.go
│       └── simulator.go
├── logs/
├── output/
├── test/
│   └── simulator_test.go
├── go.mod
├── go.sum
```

## Testing

The project includes unit tests that can be run with the following command:

```bash
make test
```

Or manually run the tests:

```bash
go test ./test/
```

## Useful Commands

### Build the project

```bash
make build
```

### Run the project

```bash
make run
```

### Clean the build and temp files

```bash
make clean
```

### Run tests

```bash
make test
```

### Format Go files

```bash
make fmt
```

### Run linting

```bash
make lint
```

### Install dependencies

```bash
make deps
```

### Run everything (lint, fmt, test, build)

```bash
make all
```

### Installing golangci-lint

To run the linter target, make sure to install golangci-lint. You can install it by running the following command:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Contributing

Contributions are welcome! Please follow these guidelines:

- Fork the repository on GitHub.
- Clone your forked repository to your local machine.
- Create a new branch for your feature or bugfix:
- Make your changes and ensure the code follows Go best practices.
- Run tests to ensure nothing is broken:
- Commit your changes with a descriptive commit message:
- Push to your fork:
- Submit a pull request to the main repository.
- Please make sure your code is well-documented and includes unit tests where appropriate.

License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
