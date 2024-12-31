package advance

import (
	"os"
	"runtime"

	"github.com/edsrzf/mmap-go"
)

type MemChunk struct {
	start int
	end   int
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
	station := ""
	temperature := 0
	prev := start
	measurements := make(map[string]*Measurements)

	for i := start; i <= end; i++ {
		if data[i] == ';' {
			station = string(data[prev:i])
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

			measurement := measurements[station]

			if measurement == nil {
				measurements[station] = &Measurements{
					Min:   temperature,
					Max:   temperature,
					Sum:   int64(temperature),
					Count: 1,
				}
			} else {
				measurement.Min = min(measurement.Min, temperature)
				measurement.Max = max(measurement.Max, temperature)
				measurement.Sum += int64(temperature)
				measurement.Count += 1
			}
			prev = i + 1
			station = ""
			temperature = 0
		}
	}
	ch <- measurements
}

func ParallelMmap(dataFilePath string) {
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
