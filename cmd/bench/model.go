package bench

// import (
// 	"encoding/json"
// )

type Command struct {
	Name string
	Cmd  []string `json:"cmd" binding:"required"`
}

type Scenario struct {
	Name   string `json:"name" binding:"required"`
	Env    map[string]string
	Before *Command
	After  *Command
	Script []*Command `json:"script" binding:"required"`
}

type Benchmark struct {
	Name        string `json:"name" binding:"required"`
	Description string
	Scenarios   []*Scenario `json:"scenarios" binding:"required"`
	Executions  int
	Alternate   bool
}
