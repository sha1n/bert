package bench

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Load(t *testing.T) {
	expected := &Benchmark{
		Name:        "test benchmark",
		Description: "",
		Executions:  10,
		Alternate:   true,
		Scenarios: []*Scenario{
			{
				Name: "scenario A",
				Env:  make(map[string]string),
				Before: &Command{
					Cmd: []string{"echo", "beforeA"},
				},
				After: &Command{
					Cmd: []string{"echo", "afterA"},
				},
				Script: []*Command{
					{
						Cmd: []string{"sleep", "2"},
					},
				},
			},
			{
				Name: "scenario B",
				Env:  make(map[string]string),
				Before: &Command{
					Cmd: []string{"echo", "beforeB"},
				},
				After: &Command{
					Cmd: []string{"echo", "afterB"},
				},
				Script: []*Command{
					{
						Cmd: []string{"sleep", "1"},
					},
				},
			},
		},
	}

	benchmark, err := Load("../../test_data/config_test_load.json")

	assert.NoError(t, err)
	assert.Equal(t, benchmark, expected)
}
