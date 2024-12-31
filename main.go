package main

import (
	"fmt"
	brute "ritikjainrj18/1BRC/solution/1_Brute"

	// final "ritikjainrj18/1BRC/solution/3_Final"
	// advance "ritikjainrj18/1BRC/solution/2_Advance"
	"time"
)

var dataFilePath = "1brc/measurements.txt"

func main() {
	start := time.Now()
	brute.Brute(dataFilePath)
	// advance.CustomMmap(dataFilePath)
	// advance.ParallelMmap(dataFilePath)
	// final.CustomMapParallelMmap(dataFilePath)
	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)
}
