# mstodo

[![Releases](https://github.com/dalyIsaac/mstodo/actions/workflows/build.yml/badge.svg)](https://github.com/dalyIsaac/mstodo/actions/workflows/build.yml) [![Pushes and Pull Requests](https://github.com/dalyIsaac/mstodo/actions/workflows/pr.yml/badge.svg)](https://github.com/dalyIsaac/mstodo/actions/workflows/pr.yml)

`mstodo` is a CLI program for using Microsoft To Do.

## Install

Download the latest release from the [Releases](https://github.com/dalyIsaac/mstodo/releases) page. Add the folder containing `mstodo` to the `PATH`.

## Usage

```txt
To see available commands, type mstodo help

Usage:
  mstodo [command]

Available Commands:
  add         Add a task
  help        Help about any command
  lists       Get a list of the task lists
  version     mstodo version
  view        View a specific list

Flags:
      --auth-timeout string   seconds to wait before giving up on authentication and exiting
      --config-dir string     config directory (default "/home/dalyisaac/.mstodo")
  -h, --help                  help for mstodo
      --port string           port for mstodo
  -t, --table-style string    the style for the table (default "Rounded")

Use "mstodo [command] --help" for more information about a command.
```

## Development

To install dependencies:

```shell
go mod vendor
```

To run your cloned repo:

```shell
go run main.go
```
