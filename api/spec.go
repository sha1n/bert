package api

// CommandSpec benchmark command execution specs
type CommandSpec struct {
	WorkingDirectory string   `json:"workingDir,omitempty" yaml:"workingDir,omitempty"`
	Cmd              []string `json:"cmd" yaml:"cmd" validate:"required"`
}

// ScenarioSpec benchmark scenario specs
type ScenarioSpec struct {
	Name             string            `json:"name" yaml:"name" validate:"required"`
	WorkingDirectory string            `json:"workingDir,omitempty" yaml:"workingDir,omitempty"`
	Env              map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	BeforeAll        *CommandSpec      `json:"beforeAll,omitempty" yaml:"beforeAll,omitempty"`
	AfterAll         *CommandSpec      `json:"afterAll,omitempty" yaml:"afterAll,omitempty"`
	BeforeEach       *CommandSpec      `json:"beforeEach,omitempty" yaml:"beforeEach,omitempty"`
	AfterEach        *CommandSpec      `json:"afterEach,omitempty" yaml:"afterEach,omitempty"`
	Command          *CommandSpec      `validate:"required,dive"`
}

// BenchmarkSpec benchmark specs top level structure
type BenchmarkSpec struct {
	Scenarios  []ScenarioSpec `json:"scenarios" yaml:"scenarios" validate:"required,min=1"`
	Executions int            `validate:"required,gte=1"`
	Alternate  bool           `json:"alternate,omitempty" yaml:"alternate,omitempty"`
	FailFast   bool           `json:"failFast,omitempty" yaml:"failFast,omitempty"`
}

// ID returns a unique identifier
func (s ScenarioSpec) ID() string {
	return s.Name
}
