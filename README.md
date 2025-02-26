# CPUs

[![PkgGoDev](https://img.shields.io/badge/-reference-blue?logo=go&logoColor=white&labelColor=505050)](https://pkg.go.dev/github.com/thediveo/cpus)
[![GitHub](https://img.shields.io/github/license/thediveo/cpus)](https://img.shields.io/github/license/thediveo/cpus)
![build and test](https://github.com/thediveo/cpus/actions/workflows/buildandtest.yaml/badge.svg?branch=master)
![Coverage](https://img.shields.io/badge/Coverage-95.2%25-brightgreen)
[![Go Report Card](https://goreportcard.com/badge/github.com/thediveo/cpus)](https://goreportcard.com/report/github.com/thediveo/cpus)

`cpus` is a small Go module for dealing with CPU lists and sets, as used
throughout several places in Linux, such as syscalls and `procfs` pseudo files.
It has been carved out from the [lxkns](https://github.com/thediveo/lxkns)
project as it is useful in applications, tools, and tests beyond lxkns.

Please refer to the [module
documentation](https://pkg.go.dev/github.com/thediveo/cpus) for usage and
details.

## Tinkering

When tinkering with the `cpus` source code base, the recommended way is a
devcontainer environment. The devcontainer specified in this repository
contains:

- `gocover` command to run all tests with coverage, updating the README coverage
  badge automatically after successful runs.
- Go package documentation is served in the background on port TCP/HTTP `6060`
  of the devcontainer.
- [`go-mod-upgrade`](https://github.com/oligot/go-mod-upgrade)
- [`goreportcard-cli`](https://github.com/gojp/goreportcard).
- [`pin-github-action`](https://github.com/mheap/pin-github-action) for
  maintaining Github Actions.

## Copyright and License

`cpus` is Copyright 2024‒25 Harald Albrecht, and licensed under the Apache
License, Version 2.0.
