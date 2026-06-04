# CPUs

[![PkgGoDev](https://img.shields.io/badge/-reference-blue?logo=go&logoColor=white&labelColor=505050)](https://pkg.go.dev/github.com/thediveo/cpus)
[![GitHub](https://img.shields.io/github/license/thediveo/cpus)](https://img.shields.io/github/license/thediveo/cpus)
![build and test](https://github.com/thediveo/cpus/actions/workflows/buildandtest.yaml/badge.svg?branch=master)
![Coverage](https://img.shields.io/badge/Coverage-96.7%25-brightgreen)
[![Go Report Card](https://goreportcard.com/badge/github.com/thediveo/cpus)](https://goreportcard.com/report/github.com/thediveo/cpus)

`cpus` is a small Go module for dealing with CPU lists and sets, as used
throughout several places in Linux, such as syscalls and inside `procfs` pseudo
files. It has been carved out from the
[lxkns](https://github.com/thediveo/lxkns) project as it is useful in
applications, tools, and tests beyond lxkns.

Please refer to the [module
documentation](https://pkg.go.dev/github.com/thediveo/cpus) for usage and
details.

## DevContainer

> [!CAUTION]
>
> Do **not** use VSCode's "~~Dev Containers: Clone Repository in Container
> Volume~~" command, as it is utterly broken by design, ignoring
> `.devcontainer/devcontainer.json`.

1. `git clone https://github.com/thediveo/cpus`
2. in VSCode: Ctrl+Shift+P, "Dev Containers: Open Workspace in Container..."
3. select `cpus.code-workspace` and off you go...

## Supported Go Versions

`clippy` supports versions of Go that are noted by the [Go release
policy](https://golang.org/doc/devel/release.html#policy), that is, major
versions _N_ and _N_-1 (where _N_ is the current major version).

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md).

## Copyright and License

`cpus` is Copyright 2024‒26 Harald Albrecht, and licensed under the Apache
License, Version 2.0.
