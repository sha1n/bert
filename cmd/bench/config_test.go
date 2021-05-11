package bench

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_LoadJson(t *testing.T) {
	test_Load(t, "../../test_data/config_test_load.json")
}

func Test_LoadYaml(t *testing.T) {
	test_Load(t, "../../test_data/config_test_load.yaml")
}

func test_Load(t *testing.T, configPath string) {
	expected := expectedBenchmarkConfig()

	benchmark, err := Load(configPath)

	assert.NoError(t, err)
	assert.Equal(t, expected, benchmark)
}

func expectedBenchmarkConfig() *Benchmark {
	return &Benchmark{
		Executions: 2,
		Alternate:  true,
		Scenarios: []*Scenario{
			{
				Name:             "scenario A",
				WorkingDirectory: "/tmp",
				Env:              map[string]string{"KEY": "value"},
				Before: &Command{
					Cmd: []string{"echo", "beforeA"},
				},
				After: &Command{
					Cmd: []string{"echo", "afterA"},
				},
				Script: []*Command{
					{
						Cmd: []string{"sleep", "1"},
					},
				},
			},
			{
				Name: "scenario B",
				Script: []*Command{
					{
						Cmd: []string{"sleep", "0"},
					},
					{
						Cmd: []string{"sleep", "1"},
					},
				},
			},
		},
	}

}
