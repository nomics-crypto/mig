# `mig`: the littlest migrator

[![GoDoc](https://godoc.org/github.com/nomics-crypto/mig?status.svg)](https://godoc.org/github.com/nomics-crypto/mig)
[![Build Status](https://travis-ci.org/nomics-crypto/mig.svg?branch=master)](https://travis-ci.org/nomics-crypto/mig)

`mig` is a minimal application migration tool. It comes with very strong conventions around how
migrations should be managed and how your application is configured. If you have a
[12 factor app](https://12factor.net/) on a cloud PAAS then `mig` should work for you.

`mig` uses raw SQL migration files so that the evolution of your code never renders an old migration obsolete, and also so that you can write heavily optimized SQL, and finally so that if anything goes wrong you can always execute a migration file directly, like this:

```
$ psql -a -f do-something.down.sql
```

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

## Development

PostgreSQL must be installed and running.

**You must create a database** for mig to use. (In these examples, I use a database called `mig`)

You must do **one** of these:

1. Set `DATABASE_URL` in the environment to a usable URL. For example: `postgres://me@localhost/mig`
2. Create `./libmig/.env` and add a line `DATABASE_URL=your db url`

This way, mig works on Heroku and Travis quite easily. Now you should be able to run:

```
$ go test ./...
```

You should also be able to try out the CLI:

```
$ go install
$ mig init
mig initialized
running 20171002164433-create-migrations
$ psql -c "SELECT version FROM migrations" -d mig
             version
----------------------------------
 20171002164433-create-migrations
(1 row)
```

p.s. I have tested mig on linux and on windows, but not on OS X. If it needs tweaks for OS X let me know.

## Mig on Heroku CI + Release

Here's how to use mig with Heroku's CI and Release Phase. In your `Gopkg.toml`:

```toml
required = ["github.com/nomics-crypto/mig"]

[metadata.heroku]
  root-package = "your/app/package/here"
  install = [ ".", "./vendor/github.com/nomics-crypto/mig"]
```

For `install` you just want to add mig on the end of what you already have. The default install is ".".

Then, in your Procfile:

```Procfile
web: your command here
release: mig up
```
