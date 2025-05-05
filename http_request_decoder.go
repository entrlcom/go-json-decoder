package json_decoder

import (
	"errors"
	"net/http"
	"strings"
)

var ErrUnsupportedMediaType = errors.New("unsupported media type")

type HTTPRequestDecoder struct {
	limit int64
}

func (x HTTPRequestDecoder) Decode(w http.ResponseWriter, r *http.Request, v any) error {
	if r.ContentLength > x.limit {
		return ErrRequestEntityTooLarge
	}

	if strings.Split(r.Header.Get("Content-Type"), ";")[0] != "application/json" {
		return ErrUnsupportedMediaType
	}

	return Decode(http.MaxBytesReader(w, r.Body, x.limit), v)
}

func NewHTTPRequestDecoder(limit int64) HTTPRequestDecoder {
	return HTTPRequestDecoder{
		limit: limit,
	}
}
