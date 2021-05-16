package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadJson(t *testing.T) {
	testLoad(t, "../test/data/spec_test_load.json")
}

func TestLoadYaml(t *testing.T) {
	testLoad(t, "../test/data/spec_test_load.yaml")
}

func TestLoadYamlWithMissingRequiredCommand(t *testing.T) {
	_, err := Load("../test/data/spec_test_load_with_missing_command.yaml")

	assert.Error(t, err)
}

func TestLoadYamlWithMissingCmdFieldOfOptionalCommand(t *testing.T) {
	_, err := Load("../test/data/spec_test_load_with_missing_cmd_of_optional_command.yaml")

	assert.Error(t, err)
}

func testLoad(t *testing.T, filePath string) {
	expected := expectedBenchmarkSpec()

	benchmark, err := Load(filePath)

	assert.NoError(t, err)
	assert.Equal(t, expected, benchmark)
}

func expectedBenchmarkSpec() *BenchmarkSpec {
	return &BenchmarkSpec{
		Executions: 10,
		Alternate:  true,
		Scenarios: []*ScenarioSpec{
			{
				Name:             "scenario A",
				WorkingDirectory: "/tmp",
				Env:              map[string]string{"KEY": "value"},
				BeforeAll: &CommandSpec{
					Cmd: []string{"echo", "setupA"},
				},
				AfterAll: &CommandSpec{
					Cmd: []string{"echo", "teardownA"},
				},
				BeforeEach: &CommandSpec{
					WorkingDirectory: "/another-path",
					Cmd:              []string{"echo", "beforeA"},
				},
				AfterEach: &CommandSpec{
					Cmd: []string{"echo", "afterA"},
				},
				Command: &CommandSpec{
					Cmd: []string{"sleep", "1"},
				},
			},
			{
				Name: "scenario B",
				Command: &CommandSpec{
					Cmd: []string{"sleep", "0"},
				},
			},
		},
	}

}
