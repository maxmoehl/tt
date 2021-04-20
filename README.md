# tt

`tt` is a cli application that can be used to track time. This README will be expanded
when I have more time.

# Configuration

The environment variable `TT_HOME_DIR` specifies where the application should look for
a configuration file and store any information it collects. If `TT_HOME_DIR` is not set
`$HOME/.tt` is used. If a file named `config.yaml` is present in the directory the
config is read form it, otherwise defaults are used. See the `config/config.go` for
more details.

# Installation

```
go install github.com/maxmoehl/tt@latest
```

# Usage

Documentation is available as part of the cli. Only calling `tt` prints out a help
section from which you can explore the different commands.