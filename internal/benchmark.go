package internal

import (
	"bufio"
	"log"
	"os"

	"github.com/sha1n/benchy/pkg"
)

type WriteReport = func(summary pkg.TracerSummary, config *Benchmark)

func Run(configFilePath string) error {
	log.Println("Starting benchy...")

	return run(configFilePath, writeReport)
}

func run(configFilePath string, writeReport WriteReport) (error error) {
	log.Println("Starting benchy...")

	config, err := loadConfig(configFilePath)
	if err != nil {
		return err
	}

	summary := Execute(config, NewContext(pkg.NewTracer()))
	writeReport(summary, config)

	return error
}

func loadConfig(configFilePath string) (rtn *Benchmark, error error) {
	log.Printf("Loading configuration file '%s'...\r\n", configFilePath)

	benchmark, err := Load(configFilePath)
	if err != nil {
		log.Println(err)
		error = err
	}

	return benchmark, error
}

func writeReport(summary pkg.TracerSummary, config *Benchmark) {
	console := bufio.NewWriter(os.Stdout)
	writer := NewTextReportWriter(console)
	writer.Write(summary, config)
}
