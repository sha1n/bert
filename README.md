[![Build Status](https://travis-ci.com/sha1n/benchy.svg?branch=master)](https://travis-ci.com/sha1n/benchy)
[![Go Report Card](https://goreportcard.com/badge/github.com/sha1n/benchy)](https://goreportcard.com/report/github.com/sha1n/benchy)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/sha1n/benchy.svg?style=flat-square)](https://github.com/sha1n/benchy/releases)

# benchy
`benchy` is a simple CLI benchmarking tool that allows you to easily compare key performance metrics of different CLI commands. It was developed very quickly for a very specific use-case I had, but it is already very useful and can easily evolve into something even better.

## Main Features
- Compare any number of commands
- Set your working directory per scenario and/or command
- Set the number of times every scenario is executed
- Set optional custom environment variables per scenario
- Set optional before/after commands for each run
- Set optional setup/teardown commands per scenario
- Choose between alternate executions and sequencial execution of the same command
- Choose between `txt` and `csv` output formats

## Usage
```bash
$ benchy --config test/data/spec_test_load.yaml

$ benchy --help   # for full options list
```

## Example Text Summary 
```bash
-------------------
 BENCHMARK SUMMARY
-------------------
       date: May 18 2021
       time: 14:44:24+03:00
  scenarios: 2
 executions: 10
  alternate: true

------------------------
 SCENARIO: 'scenario A'
------------------------
        min: 1.004s
        max: 1.008s
       mean: 1.006s
     median: 1.006s
        p90: 1.008s
     stddev: 0.001s
     errors: 0%

------------------------
 SCENARIO: 'scenario B'
------------------------
        min: 0.004s
        max: 0.004s
       mean: 0.004s
     median: 0.004s
        p90: 0.004s
     stddev: 0.000s
     errors: 0%
```

## Example Config
The config file can be either in JSON format or YAML. `benchy` assumes a file with the `yml` or `yaml` extension to be YAML, otherwise JSON is assumed. More about configuration [here](docs/configuration.md).

**Example YAML config:**
```yaml
---
alternate: true
executions: 10
scenarios:
- name: scenario A
  workingDir: "/tmp"
  env:
    KEY: value
  beforeAll:
    cmd:
    - echo
    - setupA
  afterAll:
    cmd:
    - echo
    - teardownA
  beforeEach:
    workingDir: "/another-path"
    cmd:
    - echo
    - beforeA
  afterEach:
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

**Equivalent JSON config:**
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
      "beforeAll": {
        "cmd": [
          "echo",
          "setupA"
        ]
      },
      "afterAll": {
        "cmd": [
          "echo",
          "teardownA"
        ]
      },
      "beforeEach": {
        "workingDir": "/another-path",
        "cmd": [
          "echo",
          "beforeA"
        ]
      },
      "afterEach": {
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
