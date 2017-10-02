# `mig`: the littlest migrator

[![GoDoc](https://godoc.org/github.com/nomics-crypto/mig?status.svg)](https://godoc.org/github.com/nomics-crypto/mig)
[![Build Status](https://travis-ci.org/nomics-crypto/mig.svg?branch=master)](https://travis-ci.org/nomics-crypto/mig)

`mig` is a minimal application migration tool. It comes with very strong conventions around how
migrations should be managed and how your application is configured. If you have a
[12 factor app](https://12factor.net/) on a cloud PAAS then `mig` should work for you.

## Installing

Install `mig` with Go:

```
go get -u github.com/nomics-crypto/mig
```

## Usage

`mig` gives you help:

```
$ mig
usage:
  mig <command> [arguments]

commands:
  init                       # Initialize an app with mig
  new <name-of-migration>    # Create a new migration with the given name
  up                         # Run all migrations that haven't been run
  down                       # Run all down migrations for migrations that have been run
  help                       # This usage information
```