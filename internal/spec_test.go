package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_LoadJson(t *testing.T) {
	test_Load(t, "../test_data/spec_test_load.json")
}

func Test_LoadYaml(t *testing.T) {
	test_Load(t, "../test_data/spec_test_load.yaml")
}

func test_Load(t *testing.T, filePath string) {
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
				Setup: &CommandSpec{
					Cmd: []string{"echo", "setupA"},
				},
				Teardown: &CommandSpec{
					Cmd: []string{"echo", "teardownA"},
				},
				BeforeCommand: &CommandSpec{
					WorkingDirectory: "/another-path",
					Cmd:              []string{"echo", "beforeA"},
				},
				AfterCommand: &CommandSpec{
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
