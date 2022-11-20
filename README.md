# TinyHelper

A CLI Util for TinyGo

# Prerequisites

## TinyGo

Requires TinyGo to be installed, see https://tinygo.org/getting-started/install/ for the installation guide.

## Direnv

Requires direnv to be installed, see https://direnv.net/docs/installation.html for the installation guide.

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
  flash       Flash target device
  help        Help about any command

Flags:
  -h, --help            help for tinyhelper
  -z, --helpz           tinygo help
  -t, --target string   target hardware
  -u, --update          update target

Use "tinyhelper [command] --help" for more information about a command.
```

## env command

```
tinyhelper env [-t {target}] [-u]
```

Creates a `.envrc` file with `GOROOT` and `GOFLAGS` configured for your environment and target device.

If target flag is passed with use to set target. If not passed and existing `TH_TARGET` env var is set
will use it as target. Else will prompt for target selection.

Update flag will prompt for selecting new target even if target is already set.

### Example generated .envrc

```
# TinyHelper START
export GOROOT=/home/gord/.cache/tinygo/goroot-d3a5eae46885c758dc170cc3b2ebb723ef9c0181c18efbe4c5dc3ba26d61a5ae
export GOFLAGS=-tags=cortexm,baremetal,linux,arm,rp2040,rp,pico,tinygo,math_big_pure_go,gc.conservative,scheduler.tasks,serial.usb
export TH_TARGET=pico
# TinyHelper END

```

![TinyHelper](https://github.com/gordcurrie/gifs/blob/main/tinyhelper.gif)

## flash command

```
tinyhlelper  flash [-t {target}][args] [-u]
```

Runs the TinyGo flash command, will accept and any valid arguments for TinyGo flash. Will use preconfigured target if set or prompt for target if not set.

Update flag will prompt for selecting new target even if target is already set.


# Development

# dev mode

If running via cli via `go run main.go env` TinyHelper will detect it is in dev mode and will output to `.envrc.temp` to prevent overwriting of Go environment.

# Configure `.envrc`

From the root the project run, add `export TH_TARGET=pico` to your existing `.envrc` file or create a
new file with the following command from the root of the project.

```
echo "export TH_TARGET=pico" > .envrc
```
