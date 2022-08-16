![CI](https://github.com/akatsuki-members/credit-crypto/actions/workflows/pubsub.yaml/badge.svg?branch=dev) ![GitHub release](https://img.shields.io/github/release/akatsuki-members/credit-crypto/all.svg?style=plastic)

# pubsub

This library defines the logic for publishing and subscribing to an event bus.

## How to build?

this is a library, we are not going to build it.

## How to test?

there are two options to test this library.

* running go tool

```sh
go test -race ./...
```

* running make

```sh
make test
```

## How to check with linter?

* running local using docker

```sh
make lint-docker 
```

* running local using docker and verbose

```sh
make lint-docker-verbose
```

## How to use?

## Known issues with linter

1.  File is not `gci`-ed with --skip-generated -s standard,default (gci)

if you are doing static analysis with golangci-lint and you see this issue:

```log
File is not `gci`-ed with --skip-generated -s standard,default (gci)
```

you can fix it installing [gci](https://github.com/daixiang0/gci)

```sh
go install github.com/daixiang0/gci@latest
```

and run this command over the file with the issue.

```sh
$ gci -w main.go
```

other source [here](https://github.com/golangci/golangci-lint/issues/1942)

2. File is not `gofumpt`-ed (gofumpt)

you can fix it installing [gofumpt](https://github.com/mvdan/gofumpt)

```sh
go install mvdan.cc/gofumpt@latest
```

and run this command over the file with the issue.

```sh
$ gofumpt -l -w subscribers/subscriber.go
```