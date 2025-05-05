package json_decoder

import (
	"bytes"
	"errors"
	"net/http"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/require"
)

type body struct {
	X string `json:"x"`
}

func TestDecode(t *testing.T) {
	t.Parallel()

	data := []byte(`{}`)

	require.NoError(t, Decode(bytes.NewReader(data), nil))
}

func TestDecode_ErrInvalidJSON(t *testing.T) {
	t.Parallel()

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		require.Error(t, Decode(iotest.ErrReader(errors.New("")), nil))
	})

	t.Run("io.EOF", func(t *testing.T) {
		t.Parallel()

		require.ErrorIs(t, Decode(http.NoBody, nil), ErrInvalidJSON)
	})

	t.Run("io.ErrUnexpectedEOF", func(t *testing.T) {
		t.Parallel()

		data := []byte(`{`)

		require.ErrorIs(t, Decode(bytes.NewReader(data), nil), ErrInvalidJSON)
	})

	t.Run("json.Decoder.More()", func(t *testing.T) {
		t.Parallel()

		data := []byte(`{}!`)

		require.ErrorIs(t, Decode(bytes.NewReader(data), nil), ErrInvalidJSON)
	})

	t.Run("json.SyntaxError", func(t *testing.T) {
		t.Parallel()

		data := []byte(`{"":}`)

		require.ErrorIs(t, Decode(bytes.NewReader(data), &body{}), ErrInvalidJSON)
	})

	t.Run("json.UnmarshalTypeError", func(t *testing.T) {
		t.Parallel()

		data := []byte(`{"x":0}`)

		require.ErrorIs(t, Decode(bytes.NewReader(data), &body{}), ErrInvalidJSON)
	})

	t.Run("json: unknown field", func(t *testing.T) {
		t.Parallel()

		data := []byte(`{"":0}`)

		require.ErrorIs(t, Decode(bytes.NewReader(data), &body{}), ErrInvalidJSON)
	})
}
