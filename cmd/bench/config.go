package bench

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func Load(path string) (*Benchmark, error) {
	var benchmark Benchmark

	jsonFile, err := os.Open(path)

	if err == nil {
		defer jsonFile.Close()

		bytes, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(bytes, &benchmark)
	}

	return &benchmark, err
}
