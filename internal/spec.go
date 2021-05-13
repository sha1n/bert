package internal

import (
	"encoding/json"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

type CommandSpec struct {
	WorkingDirectory string   `json:"workingDir" yaml:"workingDir"`
	Cmd              []string `json:"cmd" yaml:"cmd" validate:"required"`
}

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

type BenchmarkSpec struct {
	Scenarios  []*ScenarioSpec `json:"scenarios" yaml:"scenarios" validate:"required"`
	Executions int             `validate:"gte=1,required"`
	Alternate  bool
}

func (s *ScenarioSpec) Id() string {
	return s.Name
}

func Load(path string) (rtn *BenchmarkSpec, err error) {
	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		rtn, err = loadYaml(path)
	} else {
		rtn, err = loadJson(path)
	}

	if err == nil {
		v := validator.New()
		err = v.Struct(rtn)
	}

	return rtn, err
}

func loadJson(path string) (*BenchmarkSpec, error) {
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
