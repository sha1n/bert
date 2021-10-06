[![Go](https://github.com/sha1n/bert/actions/workflows/go.yml/badge.svg)](https://github.com/sha1n/bert/actions/workflows/go.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/sha1n/bert)
[![Go Report Card](https://goreportcard.com/badge/sha1n/bert)](https://goreportcard.com/report/sha1n/bert) 
[![Coverage Status](https://coveralls.io/repos/github/sha1n/bert/badge.svg)](https://coveralls.io/github/sha1n/bert?branch=master)
[![Release](https://img.shields.io/github/release/sha1n/bert.svg?style=flat-square)](https://github.com/sha1n/bert/releases)
[![Release Drafter](https://github.com/sha1n/bert/actions/workflows/release-drafter.yml/badge.svg)](https://github.com/sha1n/bert/actions/workflows/release-drafter.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# Bert

- [Bert](#bert)
  - [Overview](#overview)
  - [Installation](#installation)
    - [Install From a Homebrew Tap](#install-from-a-homebrew-tap)
    - [Download A Pre-Built Release](#download-a-pre-built-release)
    - [Build From Sources](#build-from-sources)
  - [Usage](#usage)
    - [Quick Ad-Hoc Benchmarks](#quick-ad-hoc-benchmarks)
    - [Using a Configuration File](#using-a-configuration-file)
    - [Directory Local Configuration (.bertconfig)](#directory-local-configuration-bertconfig)
  - [Reports](#reports)
    - [Report Formats](#report-formats)
    - [Accumulating Data](#accumulating-data)
    - [Labelling Data](#labelling-data)
    - [Understanding User & System Time Measurements](#understanding-user--system-time-measurements)
    - [Examples](#examples)
      - [Text Example](#text-example)
      - [JSON Example](#json-example)
      - [CSV Example](#csv-example)
      - [Markdown Example](#markdown-example)
      - [Raw CSV Example](#raw-csv-example)
  - [Output Control](#output-control)
  - [Shell Completion Scripts](#shell-completion-scripts)
  - [Alternatives](#alternatives)


<hr>
<img src="images/demo_800.gif" width="100%">

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
- Save results in `txt`, `json`, `csv`, `csv/raw`, `md` and `md/raw` formats
- Control your benchmark environment
  - Set optional working directory per scenario and/or command 
  - Set optional custom environment variables per scenario
  - Set optional global setup/teardown commands per scenario
  - Set optional before/after commands for each run
- Constant progress indication

## Installation
### Install From a Homebrew Tap
```bash
# Tap the formula repository (https://docs.brew.sh/Taps)
brew tap sha1n/tap

# Install bert
brew install bert

# Update bert
brew upgrade bert
```

### Download A Pre-Built Release
Download the appropriate binary and put it in your `PATH`.

```bash
# macOS single line example (assuming that '$HOME/.local/bin' is in your PATH):
bash -c 'VERSION="2.3.12" && TMP=$(mktemp -d) && cd $TMP && curl -sSL "https://github.com/sha1n/bert/releases/download/v${VERSION}/bert_${VERSION}_Darwin_x86_64.tar.gz" | tar -xz && mv bert ~/.local/bin && rm -rf $TMP'


# breaking it down

# specify which version you want
VERSION="2.3.12"
# create a temp directory
TMP=$(mktemp -d)
# change directory 
cd "$TMP"
# download and extract the archive in the temp dir
curl -sSL "https://github.com/sha1n/bert/releases/download/v${VERSION}/bert_${VERSION}_Darwin_x86_64.tar.gz" | tar -xz
# move the binary file to a dir in your PATH
mv bert ~/.local/bin

# optionally install completion scripts from ./completions

# delete the temp directory
rm -rf $TMP

# Once installed, you can update using the update command (this doesn't work when installing using Homebrew)
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

### Directory Local Configuration (.bertconfig)
When a file named `.bertconfig` exists in `bert`'s current directory and no other configuration method is specified, `bert` assumes that file is a benchmark configuration file and attempts to load specs from it.

## Reports
### Report Formats
There are three supported report formats, two of them support `raw` mode as follows. The formats are `txt`, `json`, `csv`, `csv/raw`, `md` and `md/raw`. 
- `txt` is the default report format. It contains stats per scenario and a header section that describes the main characteristics of the benchmark. `txt` format is primarily designed to be used in a terminal.
- `json` contains the same stats as `txt` does, minus the header section and is formatted as a JSON document. JSON is a very popular data representation format amongst programming languages and web applications particularly. This format is designed to help integrate `bert` reported data with other programs.
- `csv` contains the same stats in CSV format. It is especially useful when you want to accumulate stats from multiple benchmarks in a standard convenient format. In which case you can combine the `csv` format with `-o` and possibly `--header=false` if you want to accumulate data from separate runs in one file. 
- `csv/raw` is streaming raw trace events as CSV records and is useful if you want to load that data into a spreadsheet or other tools for further analysis.
- `md` and `md/raw` and similar to `csv` and `csv/raw` respectively, but write in Markdown table format.

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

| Timestamp            | Scenario | Samples | Labels    | Min   | Max    | Mean  | Median | Percentile 90 | StdDev  | User Time | System Time | Errors |
| -------------------- | -------- | ------- | --------- | ----- | ------ | ----- | ------ | ------------- | ------- | --------- | ----------- | ------ |
| 2021-06-20T21:01:24Z | curl     | 100     | vpn,wifi  | 3.5ms | 6.4ms  | 4.4ms | 4.3ms  | 5.0ms         | 552.2µs | 646.5µs   | 1.6ms       | 0%     |
| 2021-06-20T21:02:05Z | curl     | 100     | vpn,wired | 3.4ms | 16.0ms | 4.3ms | 4.1ms  | 4.9ms         | 1.3ms   | 634.0µs   | 1.6ms       | 0%     |
| 2021-06-20T21:02:33Z | curl     | 100     | wired     | 0.6ms | 8.1ms  | 1.3ms | 1.1ms  | 5.9ms         | 0.8ms   | 559.4µs   | 1.4ms       | 0%     |


### Understanding User & System Time Measurements
The `user` and `system` values are the calculated *mean* of measured user and system CPU time. It is important to understand that each measurement is the *sum* of the CPU times measured on all CPU cores and therefore can measure higher than perceived time measurements (min, max, mean, median, p90). The following report shows the measurements of two `go test` commands, one executed with `-p 1` which limits concurrency to `1` and the other with automatic parallelism. Notice how close the `user` and `system` metrics are and how they compare to the other metrics.

```
 BENCHMARK SUMMARY
     labels:
       date: Jun 15 2021
       time: 20:31:20
  scenarios: 2
 executions: 10
  alternate: true

---------------------------------------------------------------

   SCENARIO: [go test -p 1]
        min: 2.7s          mean: 2.8s        median: 2.8s
        max: 3.0s        stddev: 102.0ms        p90: 3.0s
       user: 2.3s        system: 1.2s        errors: 0%

---------------------------------------------------------------

   SCENARIO: [go test]
        min: 1.2s          mean: 1.3s        median: 1.2s
        max: 1.3s        stddev: 43.3ms         p90: 1.3s
       user: 2.8s        system: 1.6s        errors: 0%

---------------------------------------------------------------
```

### Examples
#### Text Example
```
 BENCHMARK SUMMARY
     labels: example-label
       date: Jun 12 2021
       time: 18:07:11
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

#### JSON Example
```json
{
	"records": [{
		"timestamp": "2021-06-16T20:13:07.946273Z",
		"name": "scenario A",
		"executions": 10,
		"labels": ["example-label"],
		"min": 1003598013,
		"max": 1008893354,
		"mean": 1006113519,
		"stddev": 1638733,
		"median": 1005970135,
		"p90": 1008442779,
		"user": 516700,
		"system": 1101100,
		"errorRate": 0
	}, {
		"timestamp": "2021-06-16T20:13:07.946273Z",
		"name": "scenario B",
		"executions": 10,
		"labels": ["example-label"],
		"min": 3244148,
		"max": 3907661,
		"mean": 3717243,
		"stddev": 190237,
		"median": 3795931,
		"p90": 3863124,
		"user": 544600,
		"system": 1188500,
		"errorRate": 0
	}]
}
```

#### CSV Example

```csv
Timestamp,Scenario,Samples,Labels,Min,Max,Mean,Median,Percentile 90,StdDev,User Time,System Time,Errors
2021-06-20T21:07:03Z,scenario A,10,example-label,1003901921,1007724021,1005953076,1006408481,1007591335,1316896,444700,1182600,0%
2021-06-20T21:07:03Z,scenario B,10,example-label,3256932,6759926,3939343,3626742,3947804,963627,469500,1325200,0%
```

#### Markdown Example
```
| Timestamp            | Scenario   | Samples | Labels        | Min   | Max   | Mean  | Median | Percentile 90 | StdDev  | User Time | System Time | Errors |
| -------------------- | ---------- | ------- | ------------- | ----- | ----- | ----- | ------ | ------------- | ------- | --------- | ----------- | ------ |
| 2021-06-20T21:06:22Z | scenario A | 10      | example-label | 1.0s  | 1.0s  | 1.0s  | 1.0s   | 1.0s          | 1.8ms   | 459.5µs   | 1.2ms       | 0%     |
| 2021-06-20T21:06:22Z | scenario B | 10      | example-label | 3.1ms | 4.5ms | 3.4ms | 3.2ms  | 3.7ms         | 408.8µs | 439.0µs   | 1.1ms       | 0%     |
```

#### Raw CSV Example
```csv
Timestamp,Scenario,Labels,Duration,User Time,System Time,Error
2021-06-20T21:07:25Z,scenario A,example-label,1003657149,419000,1034000,false
2021-06-20T21:07:25Z,scenario B,example-label,3234965,407000,1018000,false
2021-06-20T21:07:26Z,scenario A,example-label,1007011109,462000,1066000,false
2021-06-20T21:07:26Z,scenario B,example-label,3627639,438000,1149000,false
2021-06-20T21:07:27Z,scenario A,example-label,1004486684,480000,1057000,false
2021-06-20T21:07:27Z,scenario B,example-label,3067241,407000,942000,false
2021-06-20T21:07:28Z,scenario A,example-label,1004747640,424000,952000,false
2021-06-20T21:07:28Z,scenario B,example-label,3108394,418000,973000,false
2021-06-20T21:07:29Z,scenario A,example-label,1003754833,452000,1072000,false
2021-06-20T21:07:29Z,scenario B,example-label,3273611,403000,1047000,false
2021-06-20T21:07:30Z,scenario A,example-label,1007069722,535000,1130000,false
2021-06-20T21:07:30Z,scenario B,example-label,3110162,431000,1081000,false
2021-06-20T21:07:31Z,scenario A,example-label,1003373155,421000,959000,false
2021-06-20T21:07:31Z,scenario B,example-label,3299774,394000,987000,false
2021-06-20T21:07:32Z,scenario A,example-label,1005966685,401000,971000,false
2021-06-20T21:07:32Z,scenario B,example-label,3246400,391000,1046000,false
2021-06-20T21:07:33Z,scenario A,example-label,1004069658,462000,1074000,false
2021-06-20T21:07:33Z,scenario B,example-label,3275813,418000,1003000,false
2021-06-20T21:07:34Z,scenario A,example-label,1006268385,447000,1001000,false
2021-06-20T21:07:34Z,scenario B,example-label,3444278,401000,1054000,false
```


## Output Control
By default `bert` logs informative messages to standard err and report data to standard out (if no output file is specified). 
However, there are several ways you can control what is logged and in what level of details.

- `--pipe-stdout` and `--pipe-stderr` - pipe the standard out and err of executed benchmark commands respectively, to standard err.
- `--silent` or `-s` - sets the logging level to the lowest level possible, which includes only fatal errors. That is a softer version of `2>/dev/null` and should be preferred in general.
- `--debug` or `-d` - sets the logging level to the highest possible level, for troubleshooting.

## Shell Completion Scripts
`bert` comes with completion scripts for `zsh`, `bash`, `fish` and `PowerShell`. When installed with [brew](#install-from-a-homebrew-tap) completions scripts are automatically installed to the appropriate location, otherwise the scripts can be found in the tar-ball version of the released binaries.

## Alternatives
Before developing `bert` I looked into the following tools. Both target similar use-cases, but with different focus.
- [hyperfine](https://github.com/sharkdp/hyperfine) 
- [bench](https://github.com/Gabriel439/bench) 
