package specs

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"testing"

	"github.com/sha1n/bert/api"
	"github.com/sha1n/gommons/pkg/test"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestLoadJson(t *testing.T) {
	testLoad(t, "../../test/data/spec_test_load.json")
}

func TestLoadYaml(t *testing.T) {
	testLoad(t, "../../test/data/spec_test_load.yaml")
}

func TestLoadInvalidMissingRequiredYaml(t *testing.T) {
	_, err := LoadSpec("../../test/data/spec_test_load_invalid_missing_required.yaml")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid configuration")
}

func TestLoadInvalidTypeYaml(t *testing.T) {
	_, err := LoadSpec("../../test/data/spec_test_load_invalid_type.yaml")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot unmarshal")
}

func TestSaveYaml(t *testing.T) {
	filePath := path.Join(os.TempDir(), "TestSaveYaml.yml")
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	assert.NoError(t, err)

	expectedSpec := expectedBenchmarkSpec()

	assert.NoError(t, SaveSpec(expectedSpec, f))

	actualSpec, err := loadYaml(t, filePath)
	assert.NoError(t, err)
	assert.Equal(t, expectedSpec, actualSpec)
}

func TestSaveInvalidYaml(t *testing.T) {
	filePath := path.Join(os.TempDir(), "TestSaveYaml.yml")
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	assert.NoError(t, err)

	assert.Error(t, SaveSpec(api.BenchmarkSpec{}, f))
}

func TestSaveYamlClosesFile(t *testing.T) {
	filePath := path.Join(os.TempDir(), "TestSaveYamlClosesFile.yml")
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	assert.NoError(t, err)

	assert.NoError(t, SaveSpec(expectedBenchmarkSpec(), f))
	err = f.Close()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already closed")
}

func TestLoadYamlWithMissingRequiredCommand(t *testing.T) {
	_, err := LoadSpec("../test/data/spec_test_load_with_missing_command.yaml")

	assert.Error(t, err)
}

func TestLoadYamlWithMissingCmdFieldOfOptionalCommand(t *testing.T) {
	_, err := LoadSpec("../test/data/spec_test_load_with_missing_cmd_of_optional_command.yaml")

	assert.Error(t, err)
}

func TestLoadSpecFromYamlData(t *testing.T) {
	exec := test.RandomUint()
	scenarioName := test.RandomString()
	command := test.RandomString()

	example := fmt.Sprintf(`executions: %d
scenarios:
- name: %s
  command:
    cmd:
    - %s
`, exec, scenarioName, command)

	expected := api.BenchmarkSpec{
		Executions: int(exec),
		Scenarios: []api.ScenarioSpec{
			{
				Name:    scenarioName,
				Command: &api.CommandSpec{Cmd: []string{command}},
			},
		},
	}
	actual, err := LoadSpecFromYamlData([]byte(example))

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

}

func TestLoadSpecFromYamlDataInvalid(t *testing.T) {
	example := `executions: text
scenarios:
- name: test
  command:
    cmd:
    - test
`

	_, err := LoadSpecFromYamlData([]byte(example))

	assert.Error(t, err)
}

func TestCreateSpecFrom(t *testing.T) {
	type args struct {
		executions int
		alternate  bool
		commands   []api.CommandSpec
	}
	tests := []struct {
		name     string
		args     args
		wantSpec api.BenchmarkSpec
		wantErr  bool
	}{
		{
			name:     "valid alternate",
			args:     args{executions: 1, alternate: true, commands: []api.CommandSpec{{Cmd: []string{"ls", "-l"}}}},
			wantSpec: api.BenchmarkSpec{Executions: 1, Alternate: true, Scenarios: []api.ScenarioSpec{{Name: "[ls -l]", Command: &api.CommandSpec{Cmd: []string{"ls", "-l"}}}}},
			wantErr:  false,
		},
		{
			name:     "valid dual command",
			args:     args{executions: 1, alternate: false, commands: []api.CommandSpec{{Cmd: []string{"ls", "-l"}}, {Cmd: []string{"ls", "-a"}}}},
			wantSpec: api.BenchmarkSpec{Executions: 1, Alternate: false, Scenarios: []api.ScenarioSpec{{Name: "[ls -l]", Command: &api.CommandSpec{Cmd: []string{"ls", "-l"}}}, {Name: "[ls -a]", Command: &api.CommandSpec{Cmd: []string{"ls", "-a"}}}}},
			wantErr:  false,
		},
		{
			name:     "non-positive executions",
			args:     args{executions: rand.Intn(10) * -1, alternate: false, commands: []api.CommandSpec{{Cmd: []string{"ls", "-l"}}}},
			wantSpec: api.BenchmarkSpec{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSpec, err := CreateSpecFrom(tt.args.executions, tt.args.alternate, tt.args.commands...)

			if tt.wantErr {
				assert.Error(t, err)
			}

			assert.Equal(t, tt.wantSpec, gotSpec)
		})
	}
}

func testLoad(t *testing.T, filePath string) {
	expected := expectedBenchmarkSpec()

	benchmark, err := LoadSpec(filePath)

	assert.NoError(t, err)
	assert.Equal(t, expected, benchmark)
}

func expectedBenchmarkSpec() api.BenchmarkSpec {
	return api.BenchmarkSpec{
		Executions: 10,
		Alternate:  true,
		Scenarios: []api.ScenarioSpec{
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
					WorkingDirectory: "~/tmp",
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

func loadYaml(t *testing.T, filePath string) (spec api.BenchmarkSpec, err error) {
	var bytes []byte
	bytes, err = os.ReadFile(filePath)
	assert.NoError(t, err)

	err = yaml.Unmarshal(bytes, &spec)

	return spec, err
}
