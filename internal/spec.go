package internal

import (
	"encoding/json"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

// CommandSpec benchmark command execution specs
type CommandSpec struct {
	WorkingDirectory string   `json:"workingDir" yaml:"workingDir"`
	Cmd              []string `json:"cmd" yaml:"cmd" validate:"required"`
}

// ScenarioSpec benchmark scenario specs
type ScenarioSpec struct {
	Name             string `json:"name" yaml:"name" validate:"required"`
	WorkingDirectory string `json:"workingDir" yaml:"workingDir"`
	Env              map[string]string
	BeforeAll        *CommandSpec `json:"beforeAll" yaml:"beforeAll"`
	AfterAll         *CommandSpec `json:"afterAll" yaml:"afterAll"`
	BeforeEach       *CommandSpec `json:"beforeEach" yaml:"beforeEach"`
	AfterEach        *CommandSpec `json:"afterEach" yaml:"afterEach"`
	Command          *CommandSpec `json:"command" yaml:"command" validate:"required"`
}

// BenchmarkSpec benchmark specs top level structure
type BenchmarkSpec struct {
	Scenarios  []*ScenarioSpec `json:"scenarios" yaml:"scenarios" validate:"required"`
	Executions int             `validate:"gte=1,required"`
	Alternate  bool
}

// ID returns a unique identifier
func (s *ScenarioSpec) ID() string {
	return s.Name
}

// Load loads benchmark specs from the specified file.
func Load(path string) (rtn *BenchmarkSpec, err error) {
	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		rtn, err = loadYaml(path)
	} else {
		rtn, err = loadJSON(path)
	}

	if err == nil {
		v := validator.New()
		err = v.Struct(rtn)
	}

	return rtn, err
}

func loadJSON(path string) (*BenchmarkSpec, error) {
	var benchmark BenchmarkSpec

	jsonFile, err := os.Open(path)

	if err == nil {
		defer jsonFile.Close()

		bytes, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(bytes, &benchmark)
	}

	return &benchmark, err
}

func loadYaml(path string) (rtn *BenchmarkSpec, err error) {
	jsonFile, err := os.Open(path)

	if err == nil {
		defer jsonFile.Close()

		bytes, _ := ioutil.ReadAll(jsonFile)
		yaml.Unmarshal(bytes, &rtn)
	}

	return rtn, err
}
