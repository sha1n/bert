package pkg

import (
	"os"
	"path"
	"testing"

	"github.com/sha1n/benchy/api"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestLoadJson(t *testing.T) {
	testLoad(t, "../test/data/spec_test_load.json")
}

func TestLoadYaml(t *testing.T) {
	testLoad(t, "../test/data/spec_test_load.yaml")
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

func loadYaml(t *testing.T, filePath string) (spec *api.BenchmarkSpec, err error) {
	var bytes []byte
	bytes, err = os.ReadFile(filePath)
	assert.NoError(t, err)

	err = yaml.Unmarshal(bytes, &spec)

	return spec, err
}
