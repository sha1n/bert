[![Build Status](https://travis-ci.com/sha1n/benchy.svg?branch=master)](https://travis-ci.com/sha1n/benchy)
[![Go Report Card](https://goreportcard.com/badge/github.com/sha1n/benchy)](https://goreportcard.com/report/github.com/sha1n/benchy)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/sha1n/benchy.svg?style=flat-square)](https://github.com/sha1n/benchy/releases)

# benchy
`benchy` is a simple CLI benchmarking tool that allows you to easily compare key performance metrics of different CLI commands. It was developed very quickly for a very specific use-case I had, but it is already very useful and can easily evolve into something even better.

## Main Features
- Compare any number of commands
- Set the working directory for every scenario
- Set the number of times every scenario is executed
- Set optional custom environment variables per scenario
- Set optional before/after commands for each run
- Choose between alternate executions and sequencial execution of the same command

## Usage
```bash
$ benchy --config test_data/config_test_load.json
```

## Example Summary 
```bash
===================
 Benchmark Summary
===================
  scenarios: 2
 executions: 10
  alternate: true

=========================
 Summary of 'scenario A'
=========================
    min (s): 1.003
    max (s): 1.008
   mean (s): 1.006
 median (s): 1.005
    p90 (s): 1.007

=========================
 Summary of 'scenario B'
=========================
    min (s): 0.003
    max (s): 0.004
   mean (s): 0.004
 median (s): 0.004
    p90 (s): 0.004
```

## Example Config
The config file can be either in JSON format or YAML. `benchy` assumes a file with the `yml` or `yaml` extension to be YAML, otherwise JSON is assumed.

**Example JSON config:**
```json
{
  "alternate": true,
  "executions": 10,
  "scenarios": [
    {
      "name": "scenario A",
      "workingDir": "/tmp",
      "env": {
        "KEY": "value"
      },
      "before": {
        "cmd": [
          "echo",
          "beforeA"
        ]
      },
      "after": {
        "cmd": [
          "echo",
          "afterA"
        ]
      },
      "command": {
        "cmd": [
          "sleep",
          "1"
        ]
      }
    },
    {
      "name": "scenario B",
      "command": {
        "cmd": [
          "sleep",
          "0"
        ]
      }
    }
  ]
}
```


**Equivalent YAML config:**
```yaml
---
alternate: true
executions: 10
scenarios:
- name: scenario A
  workingDir: "/tmp"
  env:
    KEY: value
  before:
    cmd:
    - echo
    - beforeA
  after:
    cmd:
    - echo
    - afterA
  command:
    cmd:
    - sleep
    - '1'
- name: scenario B
  command:
    cmd:
    - sleep
    - '0'
```
