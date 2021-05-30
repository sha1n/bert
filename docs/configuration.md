# Benchmark Configuration 

- [Benchmark Configuration](#benchmark-configuration)
  - [Interactive Configuration Utility](#interactive-configuration-utility)
  - [Minimal Example](#minimal-example)
  - [Full Example](#full-example)
  - [Command Configuration](#command-configuration)
  - [Alternate Execution](#alternate-execution)

## Interactive Configuration Utility
An easy way to start playing with `benchy` configuration is to simply copy an [example configuration](#full-example), start modifying things and see what happens. But if you are not a YAML type of person and prefer to do it interactively, you might find the `config` utility useful. In any case, it is recommended that you go over the examples below and familiarize yourself with the different properties, so that you can get the most out of this utility.

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
working directory (inherits current) ?:
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

## Minimal Example
```yaml
---
executions: 100           # required. number of times to execute each scenario
scenarios:                # list of scenarios
- name: a scenario        # required. unique scenario name 
  workingDir: "/tmp"      # default working directory for commands executed in the context of this scenario 
  command:                # required. the benchmarked command of this scenario - the one stats are collected for
    cmd:                  # required. command line arguments.
    - command
    - --flag
    - arg1
    - arg2
```

## Full Example

```yaml
---
alternate: true           # 'true' to alternate scenario executions. More details below. (default=false)
executions: 100           # required. number of times to execute each scenario
scenarios:                # list of scenarios
- name: scenario A        # required. unique scenario name 
  workingDir: "/tmp"      # default working directory for commands executed in the context of this scenario 
  env:                    # environment variables to be set for commands executed in the context of this scenario 
    KEY: value
  beforeAll:              # command to be executed once before any other command is executed in the context of this scenario
    cmd:                  # required. command line arguments.
    - echo
    - setupA
  afterAll:               # command to be executed once after all other commands in the context of this scenario
    cmd:                  # required. command line arguments.
    - echo
    - teardownA
  beforeEach:             # command to be executed before each execution of this scenario
    workingDir: "~/path"  # working directory only for this command
    cmd:                  # required. command line arguments.
    - echo
    - beforeA
  afterEach:              # command to be executed after each execution of this scenario
    cmd:                  # required. command line arguments.
    - echo
    - afterA
  command:                # required. the benchmarked command of this scenario - the one stats are collected for
    cmd:                  # required. command line arguments.
    - sleep
    - '1'
- name: scenario B
  command:
    cmd:
    - sleep
    - '0'
```

## Command Configuration
The following elements share the same structure: `beforeAll`, `afterAll`, `beforeEach`, `afterEach`, `command`. 

`workingDir` - the `workingDir` property can be set globally for a scenario and optionally be overridden per command. If no working directory is set the default is the directory `benchy` is executed from. `~` prefix will be expanded to the current user home directory.

**Command structure:**
```yaml
  command:                # container for commands
    workingDir: "~/path"  # optional working directory for this command
    cmd:                  # required. the actual command line arguments to run
    - ls
    - -l
```

## Alternate Execution
By default `benchy` executes scenarios in sequence and according to the number of `executions` set for your benchmark. Set the `alternate` property to `true` if you want spread the different scenarios more evenly over the time line. 
Alternate execution can be helpful when:
- your benchmark runs for a very long time and external resources tend to behave differently over time
- you want some quiet time between executions of the same scenario to allow an external resource to cool down
