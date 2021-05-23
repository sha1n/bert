[![Build Status](https://travis-ci.com/sha1n/benchy.svg?branch=master)](https://travis-ci.com/sha1n/benchy)
[![Go Report Card](https://goreportcard.com/badge/github.com/sha1n/benchy)](https://goreportcard.com/report/github.com/sha1n/benchy)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/sha1n/benchy.svg?style=flat-square)](https://github.com/sha1n/benchy/releases)
[![Release Drafter](https://github.com/sha1n/benchy/actions/workflows/release-drafter.yml/badge.svg)](https://github.com/sha1n/benchy/actions/workflows/release-drafter.yml)

# benchy
`benchy` is a CLI benchmarking tool that allows you to easily compare performance metrics of different CLI commands. I developed this tool to benchmark and compare development tools and configurations on different environment setups and machine over time. It is designed to support complex scenarios that require high level of control and consistency.


- [benchy](#benchy)
  - [Main Features](#main-features)
  - [Installing](#installing)
    - [Download a prebuilt binary](#download-a-prebuilt-binary)
    - [Build your own binary](#build-your-own-binary)
  - [Usage](#usage)
  - [Configurartion](#configurartion)
  - [Report Formats](#report-formats)

## Main Features
- Compare any number of commands
- Set your working directory per scenario and/or command
- Set the number of times every scenario is executed
- Set optional custom environment variables per scenario
- Set optional before/after commands for each run
- Set optional setup/teardown commands per scenario
- Choose between alternate executions and sequencial execution of the same command
- Choose between `txt`, `csv`, `csv/raw`, `md` and `md/raw` output formats

## Installing 
### Download a prebuilt binary
Download the appropriate binary and put it in your `PATH`.

```bash
# macOS Example (assuming that '$HOME/.local/bin' is in your PATH):
curl -sSL https://github.com/sha1n/benchy/releases/latest/download/benchy-darwin-amd64 -o "$HOME/.local/bin/benchy"
```

### Build your own binary
```bash
# macOS Example (assuming that '$HOME/.local/bin' is in your PATH):
git clone git@github.com:sha1n/benchy.git
cd benchy
make 
cp ./bin/benchy-darwin-amd64 ~/.local/bin/benchy
```

## Usage
```bash
benchy --config test/data/spec_test_load.yaml

benchy --help   # for full options list
```

## Configurartion
`benchy` reads benchmark specifications from a config file. The config file can be either in YAML or JSON. `benchy` assumes a file with the `yml` or `yaml` extension to be YAML, otherwise JSON is assumed. You may create a configuratio file manually or use the `config` command to interactively generate your configuration.

More about configuration [here](docs/configuration.md).

**Benchy Config Utility Example** 
```bash
# add '-o filename.yaml' to save generated config to a file.
$ benchy config

--------------------------------
 BENCHMARK CONFIGURATION HELPER
--------------------------------

This tool is going to help you go through a benchmark configuration definition.

* annotates required input
? annotates optional input

more here: https://github.com/sha1n/benchy/blob/master/docs/configuration.md

--------------------------------

number of executions *: 30
alternate executions (false) ?: 1
scenario name *: sleepy scenario
working directory (inherits benchy's) ?:
define custom env vars? (y/n|enter):
add setup command? (y/n|enter): y
working directory (inherits scenario) ?:
command line *: echo 'preparing bedroom'
add teardown command? (y/n|enter):
add before each command? (y/n|enter): y
working directory (inherits scenario) ?:
command line *: echo 'going to sleep'
add after each command? (y/n|enter):
benchmarked command:
working directory (inherits scenario) ?:
command line *: sleep 1
add another scenario? (y/n|enter):


Writing your configuration...

scenarios:
- name: sleepy scenario
  beforeAll:
    cmd:
    - echo
    - preparing bedroom
  beforeEach:
    cmd:
    - echo
    - going to sleep
  command:
    cmd:
    - sleep
    - "1"
executions: 30
alternate: true
```

## Report Formats
There are three supported report formats; `txt`, `csv`, `csv/raw`, `md` and `md/raw`. `txt` is the default format and is primarily designed to be used in a terminal. `csv` is especially useful when you want to accumulate stats from multiple benchmarks in a CSV file. In which case you can combine the `csv` format with `-o` and possibly `--header=false`. 
`csv/raw` is streaming raw trace events as CSV records and is useful if you want to load that data into a spreadsheet or other tools for further analysis.
`md` and `md/raw` and similar to `csv` and `csv/raw` respectively, but write in Markdown table format.

Run `benchy --help` for more details.

**TXT Example:**
```bash
-------------------
 BENCHMARK SUMMARY
-------------------
     labels: example-label
       date: May 18 2021
       time: 23:34:13+03:00
  scenarios: 2
 executions: 10
  alternate: true

------------------------
 SCENARIO: 'scenario A'
------------------------
        min: 1.004s
        max: 1.007s
       mean: 1.006s
     median: 1.006s
        p90: 1.007s
     stddev: 0.001s
     errors: 0%

------------------------
 SCENARIO: 'scenario B'
------------------------
        min: 0.003s
        max: 0.004s
       mean: 0.004s
     median: 0.004s
        p90: 0.004s
     stddev: 0.000s
     errors: 0%
```


**Equivalent CSV Example:**
```csv
Timestamp,Scenario,Labels,Min,Max,Mean,Median,Percentile 90,StdDev,Errors
2021-05-18T23:38:49+03:00,scenario A,example-label,1003508458.000,1009577781.000,1006281483.700,1006164208.500,1008256954.000,2122427.909,0
2021-05-18T23:38:49+03:00,scenario B,example-label,2953009.000,4218971.000,3818925.400,3854585.000,4048263.000,317884.931,0
```

**Equivalent Markdown Example:**
```
|Timestamp|Scenario|Samples|Labels|Min|Max|Mean|Median|Percentile 90|StdDev|Errors|
|----|----|----|----|----|----|----|----|----|----|----|
|2021-05-21T16:21:13+03:00|scenario A|10|example-label|1.004s|1.010s|1.007s|1.008s|1.008s|0.002s|0%|
|2021-05-21T16:21:13+03:00|scenario B|10|example-label|0.001s|0.005s|0.004s|0.004s|0.004s|0.001s|0%|
```

**Raw CSV Example:**
```csv
Timestamp,Scenario,Labels,Duration,Error
2021-05-21T00:58:37+03:00,scenario A,example-label,1008861268,false
2021-05-21T00:58:37+03:00,scenario B,example-label,4021420,false
2021-05-21T00:58:38+03:00,scenario A,example-label,1006453206,false
2021-05-21T00:58:38+03:00,scenario B,example-label,3753389,false
2021-05-21T00:58:39+03:00,scenario A,example-label,1004680188,false
2021-05-21T00:58:39+03:00,scenario B,example-label,3780530,false
2021-05-21T00:58:40+03:00,scenario A,example-label,1005864471,false
2021-05-21T00:58:40+03:00,scenario B,example-label,3812982,false
2021-05-21T00:58:41+03:00,scenario A,example-label,1006431680,false
2021-05-21T00:58:41+03:00,scenario B,example-label,5208588,false
2021-05-21T00:58:42+03:00,scenario A,example-label,1005159913,false
2021-05-21T00:58:42+03:00,scenario B,example-label,3708653,false
2021-05-21T00:58:43+03:00,scenario A,example-label,1006895996,false
2021-05-21T00:58:43+03:00,scenario B,example-label,3261679,false
2021-05-21T00:58:44+03:00,scenario A,example-label,1008155810,false
2021-05-21T00:58:44+03:00,scenario B,example-label,3846961,false
2021-05-21T00:58:45+03:00,scenario A,example-label,1007275165,false
2021-05-21T00:58:45+03:00,scenario B,example-label,4039325,false
2021-05-21T00:58:46+03:00,scenario A,example-label,1003687652,false
2021-05-21T00:58:46+03:00,scenario B,example-label,3981022,false

```
