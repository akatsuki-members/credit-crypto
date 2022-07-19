![CI](https://github.com/akatsuki-members/credit-crypto/actions/workflows/common-endpoints.yaml/badge.svg?branch=dev) ![GitHub release](https://img.shields.io/github/release/akatsuki-members/credit-crypto/all.svg?style=plastic)

# common-endpoints

It provides a set of common endpoint for apps and microservices.

1. /health : Check service's health verifying all of its integrations.
2. /heartbeat : Check if service is alive.
3. /info: Provides version, commit hash and application name.

## How to build?

this is a library, we are not going to build it.

## How to test?

there are two options to test this library.

1. running go tool

```sh
go test -race ./...
```

2. running make

```sh
make test
```

## How to use?

1. You can add this library to your application as explained below.

```sh
go get github.com/akatsuki-members/credit-crypto/libs/common-endpoints
```

2. Using it in your service is pretty straight forward. 

* You only need to provide a router that implements this function.

```go
func(http.ResponseWriter, *http.Request)
```

* You can add one common endpoint or all of them.

```go
import "github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers/health"
import "github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers/heartbeat"
import "github.com/akatsuki-members/credit-crypto/libs/common-endpoints/internal/handlers/info"
...
mux := http.NewServeMux()
...
infoData := endpoints.Info{
	Name:    "audit-app",
	Commit:  "963e91b",
	Version: "1.5.3",
}
...
endpoints.New(mux).WithHeartbeat().WithHealth(newHypotheticalHealthChecker()).WithInfo(infoData)
...
func newHypotheticalHealthChecker() func() health.Report {
    // add your health logic here
	return func() health.Report {
        report, healthy := checkIntegrations()
		return health.Report{
			Healthy: healthy,
			Data:    report,
		}
	}
}
...
func checkIntegrations() ([]health.Item, bool) {
    healthy := true
    report := []endpoints.Item{
		{Name: "Database", Healthy: true},
		{Name: "Cache", Healthy: true},
		{Name: "OtherService", Healthy: true},
	}
    return report, healthy
}
...
```
