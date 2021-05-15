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
	WorkingDirectory string `json:"workingDir" yaml:"workingDir" validate:"required"`
	Env              map[string]string
	BeforeAll        *CommandSpec `json:"beforeAll" yaml:"beforeAll"`
	AfterAll         *CommandSpec `json:"afterAll" yaml:"afterAll"`
	BeforeEach       *CommandSpec `json:"beforeEach" yaml:"beforeEach"`
	AfterEach        *CommandSpec `json:"afterEach" yaml:"afterEach"`
	Command          *CommandSpec `validate:"required"`
}

// BenchmarkSpec benchmark specs top level structure
type BenchmarkSpec struct {
	Scenarios  []*ScenarioSpec `json:"scenarios" yaml:"scenarios" validate:"required"`
	Executions int             `validate:"gte=1"`
	Alternate  bool
}

// ID returns a unique identifier
func (s *ScenarioSpec) ID() string {
	return s.Name
}

// Load loads benchmark specs from the specified file.
func Load(path string) (*BenchmarkSpec, error) {
	var unmarshalFn func([]byte, interface{}) error

	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		unmarshalFn = yaml.Unmarshal
	} else {
		unmarshalFn = json.Unmarshal
	}

	return load(path, unmarshalFn)
}

func load(path string, unmarshal func([]byte, interface{}) error) (spec *BenchmarkSpec, err error) {
	var jsonFile *os.File
	if jsonFile, err = os.Open(path); err == nil {
		defer jsonFile.Close()

		var bytes []byte
		if bytes, err = ioutil.ReadAll(jsonFile); err == nil {
			err = unmarshal(bytes, &spec)
		}

		if err == nil {
			v := validator.New()
			err = v.Struct(spec)
		}
	}

	return spec, err
}
