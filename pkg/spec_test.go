package pkg

import (
	"testing"

	"github.com/sha1n/benchy/api"
	"github.com/stretchr/testify/assert"
)

func TestLoadJson(t *testing.T) {
	testLoad(t, "../test/data/spec_test_load.json")
}

func TestLoadYaml(t *testing.T) {
	testLoad(t, "../test/data/spec_test_load.yaml")
}

func TestLoadYamlWithMissingRequiredCommand(t *testing.T) {
	_, err := LoadSpec("../test/data/spec_test_load_with_missing_command.yaml")

	assert.Error(t, err)
}

func TestLoadYamlWithMissingCmdFieldOfOptionalCommand(t *testing.T) {
	_, err := LoadSpec("../test/data/spec_test_load_with_missing_cmd_of_optional_command.yaml")

	assert.Error(t, err)
}

func testLoad(t *testing.T, filePath string) {
	expected := expectedBenchmarkSpec()

	benchmark, err := LoadSpec(filePath)

	assert.NoError(t, err)
	assert.Equal(t, expected, benchmark)
}

func expectedBenchmarkSpec() *api.BenchmarkSpec {
	return &api.BenchmarkSpec{
		Executions: 10,
		Alternate:  true,
		Scenarios: []*api.ScenarioSpec{
			{
				Name:             "scenario A",
				WorkingDirectory: "/tmp",
				Env:              map[string]string{"KEY": "value"},
				BeforeAll: &api.CommandSpec{
					Cmd: []string{"echo", "setupA"},
				},
				AfterAll: &api.CommandSpec{
					Cmd: []string{"echo", "teardownA"},
				},
				BeforeEach: &api.CommandSpec{
					WorkingDirectory: "/another-path",
					Cmd:              []string{"echo", "beforeA"},
				},
				AfterEach: &api.CommandSpec{
					Cmd: []string{"echo", "afterA"},
				},
				Command: &api.CommandSpec{
					Cmd: []string{"sleep", "1"},
				},
			},
			{
				Name: "scenario B",
				Command: &api.CommandSpec{
					Cmd: []string{"sleep", "0"},
				},
			},
		},
	}

}
