# Benchmark Configuration 

- [Home](README.md)
- [Benchmark Configuration](#benchmark-configuration)
  - [Interactive Configuration Utility](#interactive-configuration-utility)
  - [Starting With an Example](#starting-with-an-example)
  - [Building a Full Config File Interactively](#building-a-full-config-file-interactively)
  - [Command Configuration Structure](#command-configuration-structure)
  - [Alternate Execution](#alternate-execution)

## Interactive Configuration Utility
An easy way to start playing with `bert` configuration is to simply use an [example](#starting-with-an-example), start modifying things and see what happens. But if you are not a YAML type of person and prefer to do it interactively, you might find the [interactive config utility](#building-a-full-config-file-interactively). In any case, it is recommended that you go over the examples below and familiarize yourself with the different properties, so that you can get the most out of this utility.

## Starting With an Example
bert can generate a documented YAML example configuration for you to help you get started with a by-example approach. Here is how you do it.

```bash
# Writing a configuration example to the console
$ bert config --example

# Writing a configuration example to a file and editing it immediately using vi
$ bert config --example -o bert.yml && vi bert.yml
```
**Here is what it looks like**
```
alternate: true           # 'true' to alternate scenario executions. More details below. (default=false)
executions: 100           # required. number of times to execute each scenario
scenarios:                # list of scenarios
- name: full scenario     # required. unique scenario name
  workingDir: "/dir"      # default working directory for commands executed in the context of this scenario
  env:                    # environment variables to be set for commands executed in the context of this scenario
    NAME: value
  beforeAll:              # command to be executed once before any other command is executed in the context of this scenario
    cmd:                  # required. command line arguments.
    - command
    - --flag
  afterAll:               # command to be executed once after all other commands in the context of this scenario
    cmd:                  # required. command line arguments.
    - command
    - --flag
  beforeEach:             # command to be executed before each execution of this scenario
    workingDir: "~/dir"   # working directory only for this command
    cmd:                  # required. command line arguments.
    - command
    - --flag
  afterEach:              # command to be executed after each execution of this scenario
    cmd:                  # required. command line arguments.
    - command
    - --flag
  command:                # required. the benchmarked command of this scenario - the one stats are collected for
    cmd:                  # required. command line arguments.
    - benchmarked-command
    - --flag
    - --arg=value
- name: minimal scenario
  command:
    cmd:
    - command
```


## Building a Full Config File Interactively
```bash
# add '-o filename.yaml' to save generated config to a file.
$ bert config

--------------------------------
 BENCHMARK CONFIGURATION HELPER
--------------------------------

This tool is going to help you go through a benchmark configuration definition.

* annotates required input
? annotates optional input

more here: https://github.com/sha1n/bert/blob/master/docs/configuration.md

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

## Command Configuration Structure
The following elements share the same structure: `beforeAll`, `afterAll`, `beforeEach`, `afterEach`, `command`. 

`workingDir` - the `workingDir` property can be set globally for a scenario and optionally be overridden per command. If no working directory is set the default is the directory `bert` is executed from. `~` prefix will be expanded to the current user home directory.

**Command structure:**
```yaml
  command:                # container for commands
    workingDir: "~/path"  # optional working directory for this command
    cmd:                  # required. the actual command line arguments to run
    - ls
    - -l
```

## Alternate Execution
By default `bert` executes scenarios in sequence and according to the number of `executions` set for your benchmark. Set the `alternate` property to `true` if you want spread the different scenarios more evenly over the time line. 
Alternate execution can be helpful when:
- your benchmark runs for a very long time and external resources tend to behave differently over time
- you want some quiet time between executions of the same scenario to allow an external resource to cool down
