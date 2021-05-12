package internal

import (
	"bufio"
	"log"
	"os"

	"github.com/sha1n/benchy/pkg"
)

func Run(configFilePath string) {
	log.Println("Starting benchy...")

	config := loadConfig(configFilePath)
	summary := Execute(config, NewContext(pkg.NewTracer()))
	writeReport(summary, config)
}

func loadConfig(configFilePath string) *Benchmark {
	log.Printf("Loading configuration file '%s'...\r\n", configFilePath)

	benchmark, err := Load(configFilePath)
	if err != nil {
		log.Println(err)
	}

	return benchmark
}

func writeReport(summary pkg.TracerSummary, config *Benchmark) {
	console := bufio.NewWriter(os.Stdout)
	writer := NewTextReportWriter(console)
	writer.Write(summary, config)
}
