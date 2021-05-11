package bench

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Load(t *testing.T) {
	expected := &Benchmark{
		Name:        "test benchmark",
		Description: "",
		Executions:  2,
		Alternate:   true,
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

	benchmark, err := Load("../../test_data/config_test_load.json")

	assert.NoError(t, err)
	assert.Equal(t, expected, benchmark)
}
