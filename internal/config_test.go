package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_LoadJson(t *testing.T) {
	test_Load(t, "../test_data/config_test_load.json")
}

func Test_LoadYaml(t *testing.T) {
	test_Load(t, "../test_data/config_test_load.yaml")
}

func test_Load(t *testing.T, configPath string) {
	expected := expectedBenchmarkConfig()

	benchmark, err := Load(configPath)

	assert.NoError(t, err)
	assert.Equal(t, expected, benchmark)
}

func expectedBenchmarkConfig() *Benchmark {
	return &Benchmark{
		Executions: 10,
		Alternate:  true,
		Scenarios: []*Scenario{
			{
				Name:             "scenario A",
				WorkingDirectory: "/tmp",
				Env:              map[string]string{"KEY": "value"},
				Setup: &Command{
					Cmd: []string{"echo", "setupA"},
				},
				Teardown: &Command{
					Cmd: []string{"echo", "teardownA"},
				},
				BeforeCommand: &Command{
					WorkingDirectory: "/another-path",
					Cmd:              []string{"echo", "beforeA"},
				},
				AfterCommand: &Command{
					Cmd: []string{"echo", "afterA"},
				},
				Command: &Command{
					Cmd: []string{"sleep", "1"},
				},
			},
			{
				Name: "scenario B",
				Command: &Command{
					Cmd: []string{"sleep", "0"},
				},
			},
		},
	}

}
