package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
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
	Command          *CommandSpec `validate:"required,dive"`
}

// BenchmarkSpec benchmark specs top level structure
type BenchmarkSpec struct {
	Scenarios  []*ScenarioSpec `json:"scenarios" yaml:"scenarios" validate:"required,dive"`
	Executions int             `validate:"required,gte=1"`
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
			err = validate(spec)
		}
	}

	return spec, err
}

func validate(spec *BenchmarkSpec) (err error) {
	v := validator.New()
	return v.Struct(spec)
}
