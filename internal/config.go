package internal

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

type Command struct {
	WorkingDirectory string   `json:"workingDir" yaml:"workingDir"`
	Cmd              []string `json:"cmd" yaml:"cmd" binding:"required"`
}

type Scenario struct {
	Name             string `json:"name" yaml:"name" binding:"required"`
	WorkingDirectory string `json:"workingDir" yaml:"workingDir"`
	Env              map[string]string
	Before           *Command
	After            *Command
	Command          *Command `json:"command" yaml:"command" binding:"required"`
}

type Benchmark struct {
	Scenarios  []*Scenario `json:"scenarios" yaml:"scenarios" binding:"required"`
	Executions int
	Alternate  bool
}

func (s *Scenario) Id() string {
	return s.Name
}

func Load(path string) (*Benchmark, error) {
	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		return load_yaml(path)
	} else {
		return load_json(path)
	}
}

func load_json(path string) (*Benchmark, error) {
	var benchmark Benchmark

	jsonFile, err := os.Open(path)

	if err == nil {
		defer jsonFile.Close()

		bytes, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(bytes, &benchmark)
	}

	return &benchmark, err
}

func load_yaml(path string) (*Benchmark, error) {
	var benchmark Benchmark

	jsonFile, err := os.Open(path)

	if err == nil {
		defer jsonFile.Close()

		bytes, _ := ioutil.ReadAll(jsonFile)
		yaml.Unmarshal(bytes, &benchmark)
	}

	return &benchmark, err
}
