package bench

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Command struct {
	Cmd []string `json:"cmd" binding:"required"`
}

type Scenario struct {
	Name             string `json:"name" binding:"required"`
	WorkingDirectory string `json:"workingDir"`
	Env              map[string]string
	Before           *Command
	After            *Command
	Script           []*Command `json:"script" binding:"required"`
}

type Benchmark struct {
	Name        string `json:"name" binding:"required"`
	Description string
	Scenarios   []*Scenario `json:"scenarios" binding:"required"`
	Executions  int
	Alternate   bool
}

func (s *Scenario) Id() string {
	return s.Name
}

func Load(path string) (*Benchmark, error) {
	var benchmark Benchmark

	jsonFile, err := os.Open(path)

	if err == nil {
		defer jsonFile.Close()

		bytes, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(bytes, &benchmark)
	}

	return &benchmark, err
}
