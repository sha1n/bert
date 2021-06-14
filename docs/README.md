[![Go](https://github.com/sha1n/bert/actions/workflows/go.yml/badge.svg)](https://github.com/sha1n/bert/actions/workflows/go.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/sha1n/bert)
[![Go Report Card](https://goreportcard.com/badge/sha1n/bert)](https://goreportcard.com/report/sha1n/bert) 
[![Coverage Status](https://coveralls.io/repos/github/sha1n/bert/badge.svg)](https://coveralls.io/github/sha1n/bert?branch=master)
[![Release](https://img.shields.io/github/release/sha1n/bert.svg?style=flat-square)](https://github.com/sha1n/bert/releases)
[![Release Drafter](https://github.com/sha1n/bert/actions/workflows/release-drafter.yml/badge.svg)](https://github.com/sha1n/bert/actions/workflows/release-drafter.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# bert 

<img src="images/demo_800.gif" width="100%">

- [bert](#bert)
  - [Overview](#overview)
  - [Installation](#installation)
    - [Download A Pre-Built Release](#download-a-pre-built-release)
    - [Build From Sources](#build-from-sources)
  - [Usage](#usage)
    - [Quick Ad-Hoc Benchmarks](#quick-ad-hoc-benchmarks)
    - [Using a Configuration File](#using-a-configuration-file)
  - [Reports](#reports)
    - [Report Formats](#report-formats)
    - [Accumulating Data](#accumulating-data)
    - [Labelling Data](#labelling-data)
    - [Examples](#examples)
      - [Text Example](#text-example)
      - [CSV Example](#csv-example)
      - [Markdown Example](#markdown-example)
      - [Raw CSV Example](#raw-csv-example)
  - [Output Control](#output-control)
  - [Alternatives](#alternatives)


## Overview
`bert` is a fully-featured CLI benchmarking tool that can handle anything from the simplest ad-hoc A/B command benchmarks to multi-command scenarios with custom environment variables, working directories and more. bert can report results in several [formats](#report-formats) and forms. Reports from different runs can be marked with labels and accumulated into the same report file for later analysis. This can come handy when you want to compare different environment factors like wired network and WiFi, different disks, different software versions etc.

**Key Features**
- Benchmark any number of commands
- Perceived time measurements alongside user and system CPU time measurements
- Run quick ad-hoc benchmarks or use config files to unlock all the features
- Rerun the exact same benchmark on different machines or environments using config files
- Accumulate results for different runs and compare them later
- Set the number of times every scenario is executed
- Choose between alternate executions and sequential execution of the same command
- Save results in `txt`, `csv`, `csv/raw`, `md` and `md/raw` formats
- Control your benchmark environment
  - Set optional working directory per scenario and/or command 
  - Set optional custom environment variables per scenario
  - Set optional global setup/teardown commands per scenario
  - Set optional before/after commands for each run
- Constant progress indication

## Installation
### Download A Pre-Built Release
Download the appropriate binary and put it in your `PATH`.

```bash
# macOS Example (assuming that '$HOME/.local/bin' is in your PATH):
curl -sSL https://github.com/sha1n/bert/releases/latest/download/bert-darwin-amd64 -o "$HOME/.local/bin/bert"

# once you have it, you can update using the update command
bert update
```

### Build From Sources
If you are a Go developer or have the tools required to build Go programs, you should be able to do so by following these commands.
```bash
# macOS Example (assuming that '$HOME/.local/bin' is in your PATH):
git clone git@github.com:sha1n/bert.git
cd bert
make 
cp ./bin/bert-darwin-amd64 ~/.local/bin/bert
```

## Usage
### Quick Ad-Hoc Benchmarks
Use this form if you want to quickly measure the execution time of a command or a set of commands.

```bash
# One command 
bert 'command -opt' --executions 100

# Multiple commands
bert 'command -optA' 'command -optB' 'anotherCommand' --executions 100
```

### Using a Configuration File
In order to gain full control over benchmark configuration `bert` uses a configuration file. The configuration file can be either in YAML or JSON format. `bert` treats files with the `.json` extension as JSON, otherwise it assumes YAML. You may create a configuration file manually or use the `config` command to interactively generate your configuration.

**Why use a config files?**

- unlock advanced features such as alternate execution, custom environment variables, working directories and setup commands per scenario
- easily share elaborate benchmark configurations, store them in VCS and reuse them on different environments and over time

More about configuration [here](configuration.md).

```bash
bert --config benchmark-config.yml

# Equivalent shorthand version of the above
bert -c benchmark-config.yml
```

## Reports
### Report Formats
There are three supported report formats, two of them support `raw` mode as follows. The formats are `txt`, `csv`, `csv/raw`, `md` and `md/raw`. `txt` is the default format and is primarily designed to be used in a terminal. `csv` is especially useful when you want to accumulate stats from multiple benchmarks in a standard convenient format. In which case you can combine the `csv` format with `-o` and possibly `--header=false` if you want to accumulate data from separate runs in one file. 
`csv/raw` is streaming raw trace events as CSV records and is useful if you want to load that data into a spreadsheet or other tools for further analysis.
`md` and `md/raw` and similar to `csv` and `csv/raw` respectively, but write in Markdown table format.

**Selecting Report Format:**
```bash
# The following command will generate a report in CSV format and save it into a file 
# named 'benchmark-report.csv' in the current directory.
bert --config benchmark-config.yml --format csv --out-file benchmark-report.csv

# Here is an equivalent command that uses shorthand flag names.
bert -c benchmark-config.yml -f csv -o benchmark-report.csv
```

### Accumulating Data
When an output file is specified, `bert` *appends* data to the specified report file. If you are using one of the tabular report formats and want to accumulate data from different runs into the same report, you can specify `--headers=false` starting from the second run, to indicate that you don't want table headers.

### Labelling Data
Sometimes what you really want to measure is the impact of environmental changes on your commands and not the command themselves. In such cases, it is sometimes easier to run the exact same benchmark configuration several times, with different machine configuration. For example, WiFi vs wired network, different disks, different software versions etc. In such situations, it is helpful to label your reports in a way that allows you to easily identify each run. `bert` provides the optional `--label` or `-l` flag just for that. When specified, the provided labels will be attached to the report results of all the commands in that run.

**Example:**

| Timestamp                 | Scenario | Samples | Labels    | Min   | Max    | Mean  | Median | Percentile 90 | StdDev | Errors |
| ------------------------- | -------- | ------- | --------- | ----- | ------ | ----- | ------ | ------------- | ------ | ------ |
| 2021-06-14T19:11:00+03:00 | curl     | 100     | vpn,wifi  | 3.2ms | 11.6ms | 5.1ms | 4.7ms  | 7.0ms         | 1.5ms  | 0%     |
| 2021-06-14T19:11:09+03:00 | curl     | 100     | vpn,wired | 3.7ms | 10.8ms | 5.2ms | 4.8ms  | 8.0ms         | 1.5ms  | 0%     |
| 2021-06-14T19:12:37+03:00 | curl     | 100     | wired     | 0.6ms | 8.1ms  | 1.3ms | 1.1ms  | 5.9ms         | 0.8ms  | 0%     |


### Examples
#### Text Example
```
 BENCHMARK SUMMARY
     labels: example-label
       date: Jun 12 2021
       time: 18:07:11+03:00
  scenarios: 2
 executions: 10
  alternate: true

---------------------------------------------------------------

   SCENARIO: scenario A
        min: 1.0s          mean: 1.0s        median: 1.0s
        max: 1.0s        stddev: 1.5ms          p90: 1.0s
       user: 538.5µs     system: 1.2ms       errors: 0%

---------------------------------------------------------------

   SCENARIO: scenario B
        min: 3.4ms         mean: 3.7ms       median: 3.6ms
        max: 4.3ms       stddev: 243.9µs        p90: 3.8ms
       user: 539.9µs     system: 1.2ms       errors: 0%

---------------------------------------------------------------
```

#### CSV Example

```csv
Timestamp,Scenario,Samples,Labels,Min,Max,Mean,Median,Percentile 90,StdDev,Errors
2021-06-10T16:23:00+03:00,scenario A,10,example-label,1002473555,1006631000,1004841316,1004925820,1006234538,1263756,0%
2021-06-10T16:23:00+03:00,scenario B,10,example-label,1387363,1903073,1680485,1681815,1755130,121962,0%
```

#### Markdown Example
```
| Timestamp                 | Scenario   | Samples | Labels        | Min   | Max   | Mean  | Median | Percentile 90 | StdDev  | Errors |
| ------------------------- | ---------- | ------- | ------------- | ----- | ----- | ----- | ------ | ------------- | ------- | ------ |
| 2021-06-10T16:22:26+03:00 | scenario A | 10      | example-label | 1.0s  | 1.0s  | 1.0s  | 1.0s   | 1.0s          | 1.0ms   | 0%     |
| 2021-06-10T16:22:26+03:00 | scenario B | 10      | example-label | 1.4ms | 1.8ms | 1.6ms | 1.6ms  | 1.7ms         | 119.5µs | 0%     |
```

#### Raw CSV Example
```csv
Timestamp,Scenario,Labels,Duration,Error
2021-06-10T16:23:35+03:00,scenario A,example-label,1004649721,false
2021-06-10T16:23:35+03:00,scenario B,example-label,1435355,false
2021-06-10T16:23:36+03:00,scenario A,example-label,1005821367,false
2021-06-10T16:23:36+03:00,scenario B,example-label,1408095,false
2021-06-10T16:23:37+03:00,scenario A,example-label,1005883260,false
2021-06-10T16:23:37+03:00,scenario B,example-label,1450118,false
2021-06-10T16:23:38+03:00,scenario A,example-label,1003183088,false
2021-06-10T16:23:38+03:00,scenario B,example-label,1493344,false
2021-06-10T16:23:39+03:00,scenario A,example-label,1006013983,false
2021-06-10T16:23:39+03:00,scenario B,example-label,1809823,false
2021-06-10T16:23:40+03:00,scenario A,example-label,1005611960,false
2021-06-10T16:23:40+03:00,scenario B,example-label,1651744,false
2021-06-10T16:23:41+03:00,scenario A,example-label,1004485805,false
2021-06-10T16:23:41+03:00,scenario B,example-label,1619021,false
2021-06-10T16:23:42+03:00,scenario A,example-label,1004949755,false
2021-06-10T16:23:42+03:00,scenario B,example-label,1420214,false
2021-06-10T16:23:43+03:00,scenario A,example-label,1003371958,false
2021-06-10T16:23:43+03:00,scenario B,example-label,1747274,false
2021-06-10T16:23:44+03:00,scenario A,example-label,1004888034,false
2021-06-10T16:23:44+03:00,scenario B,example-label,1342145,false
```


## Output Control
By default `bert` logs informative messages to standard err and report data to standard out (if no output file is specified). 
However, there are several ways you can control what is logged and in what level of details.

- `--pipe-stdout` and `--pipe-stderr` - pipe the standard out and err of executed benchmark commands respectively, to standard err.
- `--silent` or `-s` - sets the logging level to the lowest level possible, which includes only fatal errors. That is a softer version of `2>/dev/null` and should be preferred in general.
- `--debug` or `-d` - sets the logging level to the highest possible level, for troubleshooting.

## Alternatives
Before developing `bert` I looked into the following tools. Both target similar use-cases, but with different focus.
- [hyperfine](https://github.com/sharkdp/hyperfine) 
- [bench](https://github.com/Gabriel439/bench) 
