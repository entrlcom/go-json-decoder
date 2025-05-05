package json_decoder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	ErrInvalidJSON           = errors.New("invalid json")
	ErrRequestEntityTooLarge = errors.New("request entity too large")
)

func Decode(r io.Reader, v any) error {
	d := json.NewDecoder(r)
	d.DisallowUnknownFields()

	if err := d.Decode(&v); err != nil {
		var (
			httpMaxBytesError      *http.MaxBytesError
			jsonSyntaxError        *json.SyntaxError
			jsonUnmarshalTypeError *json.UnmarshalTypeError
		)

		switch {
		case strings.HasPrefix(err.Error(), "json: unknown field "): // https://github.com/golang/go/issues/29035
			return fmt.Errorf("%w: %w", ErrInvalidJSON, err)
		case errors.Is(err, io.EOF):
			return fmt.Errorf("%w: %w", ErrInvalidJSON, err)
		case errors.Is(err, io.ErrUnexpectedEOF): // https://github.com/golang/go/issues/25956
			return fmt.Errorf("%w: %w", ErrInvalidJSON, err)
		case errors.As(err, &httpMaxBytesError):
			return fmt.Errorf("%w: %w", ErrRequestEntityTooLarge, err)
		case errors.As(err, &jsonSyntaxError):
			return fmt.Errorf("%w: %w", ErrInvalidJSON, err)
		case errors.As(err, &jsonUnmarshalTypeError):
			return fmt.Errorf("%w: %w", ErrInvalidJSON, err)
		default:
			return fmt.Errorf("decoder.Decode: %w", err)
		}
	}

	if d.More() {
		return ErrInvalidJSON
	}

	return nil
}
