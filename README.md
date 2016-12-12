# Horo [![GoDoc](https://godoc.org/github.com/k2wanko/horo?status.svg)](https://godoc.org/github.com/k2wanko/horo)
Horo is context friendly, Simple Web framework.

# Install 

```
go get -u github.com/k2wanko/horo
```

# Usage

```go
package main

import (
    "net/http"

    "golang.org/x/net/context"

    "github.com/k2wanko/horo"
    "github.com/k2wanko/horo/middleware"
)

func Index(c context.Context) error {
    return horo.Text(c, http.StatusOK, "Hello World!")
}

func main() {
    h := horo.New()
    h.Use(middleware.Logger(), middleware.Recover())

    h.GET("/", Index)

    h.ListenAndServe(":8080")
}
```

more [examples](https://github.com/k2wanko/horo-example)

# License

[MIT](https://github.com/k2wanko/horo/blob/master/LICENSE)