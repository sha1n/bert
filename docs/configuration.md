# Benchmark Configuration 

- [Benchmark Configuration](#benchmark-configuration)
  - [Minimal Example](#minimal-example)
  - [Full Example](#full-example)
  - [Command Configuration](#command-configuration)
  - [Alternate Execution](#alternate-execution)

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
    workingDir: "/path"   # working directory only for this command
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

`workingDir` - the `workingDir` property can be set globally for a scenario and optionally be overriden per command. If no working directory is set the default is the directory `benchy` is executed from.

**Command structure:**
```yaml
  command:                # required. the benchmarked command of this scenario - the one stats are collected for
    workingDir: "/path"   # optional working directory for this command. 
    cmd:                  # required command line arguments.
    - ls
    - -l
```

## Alternate Execution
By default `benchy` executes scenarios in sequence and according to the number of `executions` set for your benchmark. Set the `alternate` property to `true` if you want spread the different scenarios more evenly over the time line. 
Alternate execution can be helpful when:
- your benchmark runs for a very long time and external resources tend to behave differently over time
- you want some quiet time between executions of the same scenario to allow an extenral resource to cool down
