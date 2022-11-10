# TinyHelper

A CLI Util for TinyGo

# Instillation

```
go install github.com/gordcurrie/tinyhelper@latest
```

# Usage

```
Tool for helping configure tinygo

Usage:
  tinyhelper [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  env         Configures .envrc
  help        Help about any command

Flags:
  -h, --help            help for tinyhelper
  -t, --target string   target hardware

Use "tinyhelper [command] --help" for more information about a command.
```

## env command

```
tinyhelper env [-t {target}]
```

Creates a `.envrc` file with `GOROOT` and `GOFLAGS` configured for your environment and target device.

If target flag is passed with use to set target. If not passed and existing `TH_TARGET` env var is set
will use it as target. Else will prompt for target selection.

# Development

# dev mode

If running via cli via `go run main.go env` TinyHelper will detect it is in dev mode and will output to `.envrc.temp` to prevent overwriting of Go environment.

# `.envrc`

From the root the project run, add `export TH_TARGET=pico` to your existing `.envrc` file or create a
new file with the following command from the root of the project.

```
echo "export TH_TARGET=pico" > .envrc
```
