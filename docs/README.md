[![Go](https://github.com/sha1n/benchy/actions/workflows/go.yml/badge.svg)](https://github.com/sha1n/benchy/actions/workflows/go.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/sha1n/benchy)
[![Go Report Card](https://goreportcard.com/badge/sha1n/benchy)](https://goreportcard.com/report/sha1n/benchy) 
[![Coverage Status](https://coveralls.io/repos/github/sha1n/benchy/badge.svg)](https://coveralls.io/github/sha1n/benchy?branch=master)
[![Release](https://img.shields.io/github/release/sha1n/benchy.svg?style=flat-square)](https://github.com/sha1n/benchy/releases)
[![Release Drafter](https://github.com/sha1n/benchy/actions/workflows/release-drafter.yml/badge.svg)](https://github.com/sha1n/benchy/actions/workflows/release-drafter.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


# Benchy
`benchy` is a CLI benchmarking tool that allows you to easily compare performance metrics of different CLI commands. I developed this tool to benchmark and compare development tools and configurations on different environment setups and machine over time. It is designed to support complex scenarios that require high level of control and consistency.

<img src="images/demo_800.gif" width="100%">


- [Benchy](#benchy)
  - [Overview](#overview)
  - [Main Features](#main-features)
  - [Installation](#installation)
    - [Download A Pre-Built Release](#download-a-pre-built-release)
    - [Build From Sources](#build-from-sources)
  - [Usage](#usage)
    - [Quick Ad-Hoc Benchmarks](#quick-ad-hoc-benchmarks)
    - [Using a Configuration File](#using-a-configuration-file)
  - [Report Formats](#report-formats)
    - [Text Example](#text-example)
    - [CSV Example](#csv-example)
    - [Markdown Example](#markdown-example)
    - [Raw CSV Example](#raw-csv-example)
  - [Output Control](#output-control)
  - [Alternatives](#alternatives)


## Overview
`benchy` is designed with focus on benchmark environment control and flexibility in mind. It was originally built to:
- Benchmark complex, relatively long running commands such as build and test commands used on software development environments.
- Benchmark the exact same set of command scenarios on different machines or environments in order to compare them later.
- Collect raw metrics and use external analysis tools to process them.

## Main Features
- Benchmark any number of commands
- Perceived time measurements and low level user/system CPU time measurement
- Rerun the exact same benchmark again and again on different machines or environments, accumulate results and compare them later
- Set the number of times every scenario is executed
- Choose between alternate executions and sequential execution of the same command
- Save results in `txt`, `csv`, `csv/raw`, `md` and `md/raw` formats
- Control your benchmark environment
  - Set your working directory per scenario and/or command 
  - Set optional custom environment variables per scenario
  - Set optional setup/teardown commands per scenario
  - Set optional before/after commands for each run
- Constant progress indication

## Installation
### Download A Pre-Built Release
Download the appropriate binary and put it in your `PATH`.

```bash
# macOS Example (assuming that '$HOME/.local/bin' is in your PATH):
curl -sSL https://github.com/sha1n/benchy/releases/latest/download/benchy-darwin-amd64 -o "$HOME/.local/bin/benchy"

# once you have it, you can update using the update command
benchy update
```

### Build From Sources
```bash
# macOS Example (assuming that '$HOME/.local/bin' is in your PATH):
git clone git@github.com:sha1n/benchy.git
cd benchy
make 
cp ./bin/benchy-darwin-amd64 ~/.local/bin/benchy
```

## Usage
### Quick Ad-Hoc Benchmarks
Use this form if you want to quickly measure the execution time of a command or a set of commands.

```bash
# One command 
benchy 'command -opt' --executions 100

# Multiple commands
benchy 'command -optA' 'command -optB' 'anotherCommand' --executions 100
```

### Using a Configuration File
In order to gain full control over benchmark configuration `benchy` uses a configuration file. The configuration file can be either in YAML or JSON format. `benchy` treats files with the `.json` extension as JSON, otherwise it assumes YAML. You may create a configuration file manually or use the `config` command to interactively generate your configuration.

**Why use a config files?**

- unlock advanced features such as alternate execution, custom environment variables, working directories and setup commands per scenario
- easily share elaborate benchmark configurations, store them in VCS and reuse them on different environments and over time

More about configuration [here](configuration.md).

```bash
benchy --config benchmark-config.yml

# Equivalent shorthand version of the above
benchy -c benchmark-config.yml
```

## Report Formats
There are three supported report formats, two of them support `raw` mode as follows. The formats are `txt`, `csv`, `csv/raw`, `md` and `md/raw`. `txt` is the default format and is primarily designed to be used in a terminal. `csv` is especially useful when you want to accumulate stats from multiple benchmarks in a standard convenient format. In which case you can combine the `csv` format with `-o` and possibly `--header=false` if you want to accumulate data from separate runs in one file. 
`csv/raw` is streaming raw trace events as CSV records and is useful if you want to load that data into a spreadsheet or other tools for further analysis.
`md` and `md/raw` and similar to `csv` and `csv/raw` respectively, but write in Markdown table format.

**Selecting Report Format:**
```bash
# The following command will generate a report in CSV format and save it into a file 
# named 'benchmark-report.csv' in the current directory.
benchy --config benchmark-config.yml --format csv --out-file benchmark-report.csv

# Here is an equivalent command that uses shorthand flag names.
benchy -c benchmark-config.yml -f csv -o benchmark-report.csv
```

### Text Example
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

### CSV Example

```csv
Timestamp,Scenario,Samples,Labels,Min,Max,Mean,Median,Percentile 90,StdDev,Errors
2021-06-10T16:23:00+03:00,scenario A,10,example-label,1002473555,1006631000,1004841316,1004925820,1006234538,1263756,0%
2021-06-10T16:23:00+03:00,scenario B,10,example-label,1387363,1903073,1680485,1681815,1755130,121962,0%
```

### Markdown Example
```
| Timestamp                 | Scenario   | Samples | Labels        | Min   | Max   | Mean  | Median | Percentile 90 | StdDev  | Errors |
| ------------------------- | ---------- | ------- | ------------- | ----- | ----- | ----- | ------ | ------------- | ------- | ------ |
| 2021-06-10T16:22:26+03:00 | scenario A | 10      | example-label | 1.0s  | 1.0s  | 1.0s  | 1.0s   | 1.0s          | 1.0ms   | 0%     |
| 2021-06-10T16:22:26+03:00 | scenario B | 10      | example-label | 1.4ms | 1.8ms | 1.6ms | 1.6ms  | 1.7ms         | 119.5µs | 0%     |
```

### Raw CSV Example
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
By default `benchy` logs informative messages to standard err and report data to standard out (if no output file is specified). 
However, there are several ways you can control what is logged and in what level of details.

- `--pipe-stdout` and `--pipe-stderr` - pipe the standard out and err of executed benchmark commands respectively, to standard err.
- `--silent` or `-s` - sets the logging level to the lowest level possible, which includes only fatal errors. That is a softer version of `2>/dev/null` and should be preferred in general.
- `--debug` or `-d` - sets the logging level to the highest possible level, for troubleshooting.

## Alternatives
Before developing `benchy` I looked into the following tools. Both target similar use-cases, but with different focus.
- [hyperfine](https://github.com/sharkdp/hyperfine) 
- [bench](https://github.com/Gabriel439/bench) 
