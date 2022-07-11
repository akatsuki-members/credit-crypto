# golog

The golog library provides a fast and simple logger dedicated to JSON output.

golog's library is a wrapper of [Zerolog](https://pkg.go.dev/github.com/rs/zerolog#) to desacople the solutions with the Zerolog API 

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
go get github.com/akatsuki-members/credit-crypto/libs/golog
```

2. Simple logging example

```
package main

import (
	"os"

	log "github.com/akatsuki-members/credit-crypto/libs/golog"
)

func main() {
	//log setup
	l, err := log.New(log.Debug, os.Stdout)
	if err != nil {
		panic(err)
	}

	l.Info("Hello, World!")
	//output:   {"level":"info","message":"Hello, World!"}

	str := log.LogParam{Key: "str", Value: "valueStr"}
	b := log.LogParam{Key: "keyBool", Value: true}
	f := log.LogParam{Key: "keyFloat", Value: 1.23}
	i := log.LogParam{Key: "keyInt", Value: -123}
	u := log.LogParam{Key: "keyUint", Value: uint(123)}

	l.Info("With a lot of params", str, b, f, i, u)
	//output    {"level":"info","str":"valueStr","keyBool":true,"keyFloat":1.23,"keyInt":-123,"keyUint":123,"message":"With a lot of params"}

	m := map[string]string{
		"m1": "v1",
		"m2": "v2",
		"m3": "v3",
	}

	c := log.LogParam{Key: "KeyMap", Value: m}

	l.Info("With maps", c)
	//output:   {"level":"info","KeyMap":{"m1":"v1","m2":"v2","m3":"v3"},"message":"With maps"}

}

```

* Golog allows for logging at the following levels (from highest to lowest):
- Trace (value level -1)
- Debug (value level 0)
- Info  (value level 1)
- Warn  (value level 2)
- Error (value level 3)
- Fatal (value level 4)
- Panic (value level 5)

You can set the logging level at the moment of the creation of the log, in this example tje log level is Debug:
```
l, err := log.New(log.Debug, os.Stdout)
```

3. Error logging example

When an Error is sended to the logger so the output will contains the stacktrace, i.e
```
import (
	"errors"
	"os"

	log "github.com/akatsuki-members/credit-crypto/libs/golog"
)

func main() {
	//log setup
	l, err := log.New(log.Debug, os.Stdout)
	if err != nil {
		panic(err)
	}

	err = errors.New("duplicated row")
	l.Error(err, "additional msg error", log.LogParam{Key: "key", Value: "value"})

	// output    {"level":"error","key":"value","stack":[{"func":"(*Logger).Error","line":"85","source":"log.go"},{"func":"main","line":"18","source":"main.go"},{"func":"main","line":"250","source":"proc.go"},{"func":"goexit","line":"1263","source":"asm_arm64.s"}],"error":"err context: duplicated row","message":"additional msg error"}
}
```

4. Logging Fatal Messages

When fatal error is called so the golog execute an exit status 1

```
package main

import (
	"errors"
	"os"

	log "github.com/akatsuki-members/credit-crypto/libs/golog"
)

func main() {
	//log setup
	l, err := log.New(log.Debug, os.Stdout)
	if err != nil {
		panic(err)
	}

	err = errors.New("error cant open the port")
	l.Fatal(err, "fatal msg error", log.LogParam{Key: "key", Value: "value"})

    //output   {"level":"fatal","key":"value","stack":[{"func":"(*Logger).Fatal","line":"93","source":"log.go"},{"func":"main","line":"21","source":"main.go"},{"func":"main","line":"250","source":"proc.go"},{"func":"goexit","line":"1263","source":"asm_arm64.s"}],"error":"err context: error cant open the port","message":"fatal msg error"} 
     //exit status 1
}
```

4. Logging Panic Messages

When panic error is called so the golog execute an exit status 2

```
package main

import (
	"errors"
	"os"

	log "github.com/akatsuki-members/credit-crypto/libs/golog"
)

func main() {
	//log setup
	l, err := log.New(log.Debug, os.Stdout)
	if err != nil {
		panic(err)
	}

	err = errors.New("error programming defect")

	l.Panic(err, "panic msg error", log.LogParam{Key: "key", Value: "value"})
}    

Output:
{"level":"panic","key":"value","stack":[{"func":"(*Logger).Panic","line":"101","source":"log.go"},{"func":"main","line":"21","source":"main.go"},{"func":"main","line":"250","source":"proc.go"},{"func":"goexit","line":"1263","source":"asm_arm64.s"}],"error":"err context: error programming defect","message":"panic msg error"}
panic: panic msg error

goroutine 1 [running]:
...
exit status 2


```