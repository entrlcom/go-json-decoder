# JSON Decoder

## Table of Content

- [Authors](#authors)
- [Examples](#examples)

## Authors

| Name         | GitHub                                             |
|--------------|----------------------------------------------------|
| Klim Sidorov | [@entrlcom-klim](https://github.com/entrlcom-klim) |

## Examples

```go
package main

import (
	"net/http"
	"time"

	"flida.dev/unit"

	"flida.dev/json-decoder"
)

const limit = unit.B * 128 // 128 B.

type HTTPRequestBody struct {
	DateOfBirth time.Time `json:"date_of_birth"`
	Name        string    `json:"name"`
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body HTTPRequestBody

	if err := json_decoder.NewHTTPRequestDecoder(limit).Decode(w, r, &body); err != nil {
		// TODO: ...
		return
    }

	// ...
}

```
