package api

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
