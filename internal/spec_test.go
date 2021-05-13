package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadJson(t *testing.T) {
	testLoad(t, "../test_data/spec_test_load.json")
}

func TestLoadYaml(t *testing.T) {
	testLoad(t, "../test_data/spec_test_load.yaml")
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
