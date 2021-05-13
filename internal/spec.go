package internal

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

type CommandSpec struct {
	WorkingDirectory string   `json:"workingDir" yaml:"workingDir"`
	Cmd              []string `json:"cmd" yaml:"cmd" binding:"required"`
}

type ScenarioSpec struct {
	Name             string `json:"name" yaml:"name" binding:"required"`
	WorkingDirectory string `json:"workingDir" yaml:"workingDir"`
	Env              map[string]string
	Setup            *CommandSpec
	Teardown         *CommandSpec
	BeforeCommand    *CommandSpec `json:"beforeCommand" yaml:"beforeCommand" binding:"required"`
	AfterCommand     *CommandSpec `json:"aftercommand" yaml:"afterCommand" binding:"required"`
	Command          *CommandSpec `json:"command" yaml:"command" binding:"required"`
}

type BenchmarkSpec struct {
	Scenarios  []*ScenarioSpec `json:"scenarios" yaml:"scenarios" binding:"required"`
	Executions int
	Alternate  bool
}

func (s *ScenarioSpec) Id() string {
	return s.Name
}

func Load(path string) (*BenchmarkSpec, error) {
	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		return load_yaml(path)
	} else {
		return load_json(path)
	}
}

func load_json(path string) (*BenchmarkSpec, error) {
	var benchmark BenchmarkSpec

	jsonFile, err := os.Open(path)

	if err == nil {
		defer jsonFile.Close()

		bytes, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(bytes, &benchmark)
	}

	return &benchmark, err
}

func load_yaml(path string) (*BenchmarkSpec, error) {
	var benchmark BenchmarkSpec

	jsonFile, err := os.Open(path)

	if err == nil {
		defer jsonFile.Close()

		bytes, _ := ioutil.ReadAll(jsonFile)
		yaml.Unmarshal(bytes, &benchmark)
	}

	return &benchmark, err
}
