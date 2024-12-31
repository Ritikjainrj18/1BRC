package final

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"slices"

	"github.com/edsrzf/mmap-go"
)

const (
	offset64   = 14695981039346656037
	prime64    = 1099511628211
	numBuckets = 1 << 17
)

type Measurements struct {
	Min   int
	Max   int
	Sum   int64
	Count int
}

type MemChunk struct {
	start int
	end   int
}

type item struct {
	key  []byte
	stat *Measurements
}

func splitMem(mem mmap.MMap, n int) []MemChunk {
	total := len(mem)
	chunkSize := total / n
	chunks := make([]MemChunk, n)

	chunks[0].start = 0
	for i := 1; i < n; i++ {
		for j := i * chunkSize; j < i*chunkSize+50; j++ {
			if mem[j] == '\n' {
				chunks[i-1].end = j
				chunks[i].start = j + 1
				break
			}
		}
	}
	chunks[n-1].end = total - 1
	return chunks
}

func readMemChunk(ch chan map[string]*Measurements, data mmap.MMap, start int, end int) {
	temperature := 0
	prev := start
	hash := uint64(offset64)
	items := make([]item, numBuckets)

	for i := start; i <= end; i++ {
		hash ^= uint64(data[i]) // FNV-1a is XOR then *
		hash *= prime64
		if data[i] == ';' {
			station_bytes := data[prev:i]
			temperature = 0
			i += 1
			negative := false

			for data[i] != '\n' {
				ch := data[i]
				if ch == '.' {
					i += 1
					continue
				}
				if ch == '-' {
					negative = true
					i += 1
					continue
				}
				ch -= '0'
				if ch > 9 {
					panic("Invalid character")
				}
				temperature = temperature*10 + int(ch)
				i += 1
			}

			if negative {
				temperature = -temperature
			}

			hashIndex := int(hash & uint64(numBuckets-1))
			for {
				if items[hashIndex].key == nil {
					// Found empty slot, add new item.
					items[hashIndex] = item{
						key: station_bytes,
						stat: &Measurements{
							Min:   temperature,
							Max:   temperature,
							Sum:   int64(temperature),
							Count: 1,
						},
					}
					break
				}
				if bytes.Equal(items[hashIndex].key, station_bytes) {
					// Found matching slot, add to existing stats.
					s := items[hashIndex].stat
					s.Min = min(s.Min, temperature)
					s.Max = max(s.Max, temperature)
					s.Sum += int64(temperature)
					s.Count++
					break
				}
				// Slot already holds another key, try next slot (linear probe).
				hashIndex++
				if hashIndex >= numBuckets {
					hashIndex = 0
				}

			}
			prev = i + 1
			hash = uint64(offset64)
			temperature = 0
		}
	}

	measurements := make(map[string]*Measurements)
	for _, item := range items {
		if item.key == nil {
			continue
		}
		measurements[string(item.key)] = item.stat
	}
	ch <- measurements
}

func CustomMapParallelMmap(dataFilePath string) {
	maxGoroutines := min(runtime.NumCPU(), runtime.GOMAXPROCS(0))

	dataFile, err := os.Open(dataFilePath)
	if err != nil {
		panic(err)
	}
	defer dataFile.Close()

	data, err := mmap.Map(dataFile, mmap.RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer data.Unmap()

	chunks := splitMem(data, maxGoroutines)
	totals := make(map[string]*Measurements)
	measurementChan := make(chan map[string]*Measurements)

	for i := 0; i < maxGoroutines; i++ {
		go readMemChunk(measurementChan, data, chunks[i].start, chunks[i].end)
	}

	for i := 0; i < maxGoroutines; i++ {
		measurements := <-measurementChan
		for station, measurement := range measurements {
			total := totals[station]
			if total == nil {
				totals[station] = measurement
			} else {
				total.Min = min(total.Min, measurement.Min)
				total.Max = max(total.Max, measurement.Max)
				total.Sum += measurement.Sum
				total.Count += measurement.Count
			}
		}
	}

	printResults(totals)
}

func printResults(results map[string]*Measurements) {
	// sort by station name
	stationNames := make([]string, 0, len(results))
	for stationName := range results {
		stationNames = append(stationNames, stationName)
	}

	slices.Sort(stationNames)

	fmt.Printf("{")
	for idx, stationName := range stationNames {
		measurement := results[stationName]
		mean := float64(measurement.Sum/10) / float64(measurement.Count)
		max := float64(measurement.Max) / 10
		min := float64(measurement.Min) / 10
		fmt.Printf("%s=%.1f/%.1f/%.1f", stationName, min, mean, max)
		if idx < len(stationNames)-1 {
			fmt.Printf(", ")
		}
	}
	fmt.Printf("}\n")
}
