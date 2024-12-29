package brute

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type SimpleMeasurements struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

func Brute(dataFilePath string) {

	dataFile, err := os.Open(dataFilePath)
	if err != nil {
		panic(err)
	}
	defer dataFile.Close()

	measurements := make(map[string]*SimpleMeasurements)
	fileScanner := bufio.NewScanner(dataFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		rawString := fileScanner.Text()
		stationName, temperatureStr, found := strings.Cut(rawString, ";")
		if !found {
			continue
		}
		temperature, err := strconv.ParseFloat(temperatureStr, 32)
		if err != nil {
			panic(err)
		}

		measurement := measurements[stationName]
		if measurement == nil {
			measurements[stationName] = &SimpleMeasurements{
				Min:   temperature,
				Max:   temperature,
				Sum:   temperature,
				Count: 1,
			}
		} else {
			measurement.Min = min(measurement.Min, temperature)
			measurement.Max = max(measurement.Max, temperature)
			measurement.Sum += temperature
			measurement.Count += 1
		}
	}
	PrintResults(measurements)
}

func PrintResults(results map[string]*SimpleMeasurements) {
	stationNames := make([]string, 0, len(results))
	for stationName := range results {
		stationNames = append(stationNames, stationName)
	}

	slices.Sort(stationNames)
	fmt.Printf("{")
	for idx, stationName := range stationNames {
		measurement := results[stationName]
		mean := measurement.Sum / float64(measurement.Count)
		fmt.Printf("%s=%.1f/%.1f/%.1f",
			stationName, measurement.Min, mean, measurement.Max)

		if idx < len(stationNames)-1 {
			fmt.Printf(", ")
		}
	}
	fmt.Printf("}\n")
}

// real    22m2.393s
// user    2m34.303s
// sys     0m35.381s
