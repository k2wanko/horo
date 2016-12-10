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

    "github/k2wanko/horo"
)

func Index(c context.Context) error {
    return horo.Text(c, http.StatusOK, "Hello World!")
}

func main() {
    h := horo.New()

    h.GET("/", Index)

    h.ListenAndServe(":8080")
}
```

# License

[MIT](https://github.com/k2wanko/horo/blob/master/LICENSE)