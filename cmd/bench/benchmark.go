package bench

import (
	"bufio"
	"log"
	"os"
)

func Run(configFilePath string) {
	log.Println("Starting benchy...")

	config := loadConfig(configFilePath)
	summary := Execute(config, NewContext(NewTracer()))
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

func writeReport(summary TracerSummary, config *Benchmark) {
	console := bufio.NewWriter(os.Stdout)
	writer := NewTextReportWriter(console)
	writer.Write(summary, config)
}
