package main

import memorymap "ritikjainrj18/1BRC/solution/Advance"

var dataFilePath = "1brc/measurements.txt"

func main() {
	// brute.Brute(dataFilePath)
	// memorymap.CustomMmap(dataFilePath)
	memorymap.ParallelMmap(dataFilePath)
}
