package main

import (
	"fmt"
	"os"

	"github.com/sha1n/benchy/cmd/bench"
)

func main() {

	fmt.Println("Starting benchy...")

	benchmark, err := bench.Load("test_data/config_test_load.json")
	if err != nil {
		os.Exit(1)
	}

	bench.Execute(benchmark)
}
