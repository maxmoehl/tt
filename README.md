# tt

[![Go Reference](https://pkg.go.dev/badge/github.com/maxmoehl/tt.svg)](https://pkg.go.dev/github.com/maxmoehl/tt)
[![Go Report Card](https://goreportcard.com/badge/github.com/maxmoehl/tt)](https://goreportcard.com/report/github.com/maxmoehl/tt)

`tt` is a cli application that can be used to track time. This README will be expanded when I have more time.

# Configuration

The environment variable `TT_HOME_DIR` specifies where the application should look for
a configuration file and store any information it collects. If `TT_HOME_DIR` is not set
`$HOME/.tt` is used. If a file named `config.json` is present in the directory the
config is read from it, otherwise defaults are used.

# Installation

We are currently lacking automated tests therefore install this application is at your
own risk. To get the most recent version (that is somewhat manually tested) set a valid
version instead of `latest`, i.e. `v0.1.1`.

```
$ go install github.com/maxmoehl/tt@latest
```

If you already installed a version, you can check with the `--version` flag which version
you have installed.

```
$ tt --version
tt version v0.1.1
```

# Usage

Documentation is available as part of the cli. Only calling `tt` prints out a help
section from which you can explore the different commands.

