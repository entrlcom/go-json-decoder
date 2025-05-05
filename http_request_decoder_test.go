package json_decoder

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPRequestDecoder_Decode(t *testing.T) {
	t.Parallel()

	data := []byte(`{}`)

	r := newHTTPRequest(t, bytes.NewReader(data))

	w := httptest.NewRecorder()

	require.NoError(t, NewHTTPRequestDecoder(int64(len(data))).Decode(w, r, nil))
}

func TestHTTPRequestDecoder_Decode_ErrRequestEntityTooLarge(t *testing.T) {
	t.Parallel()

	t.Run("content length limit exceeded", func(t *testing.T) {
		t.Parallel()

		data := []byte(`{}`)

		r := newHTTPRequest(t, bytes.NewReader(data))

		w := httptest.NewRecorder()

		require.ErrorIs(t, NewHTTPRequestDecoder(0).Decode(w, r, nil), ErrRequestEntityTooLarge)
	})

	t.Run("http.MaxBytesError", func(t *testing.T) {
		t.Parallel()

		data := []byte(`{}`)

		r := newHTTPRequest(t, bytes.NewReader(data))
		r.ContentLength = 0
		r.Header.Set("Content-Length", "0")

		w := httptest.NewRecorder()

		require.ErrorIs(t, NewHTTPRequestDecoder(0).Decode(w, r, nil), ErrRequestEntityTooLarge)
	})
}

func TestHTTPRequestDecoder_Decode_ErrUnsupportedMediaType(t *testing.T) {
	t.Parallel()

	data := []byte(`{}`)

	r := newHTTPRequest(t, bytes.NewReader(data))
	r.Header.Del("Content-Type")

	w := httptest.NewRecorder()

	require.ErrorIs(t, NewHTTPRequestDecoder(int64(len(data))).Decode(w, r, nil), ErrUnsupportedMediaType)
}

func newHTTPRequest(t *testing.T, body io.Reader) *http.Request {
	t.Helper()

	r, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/", body)
	require.NoError(t, err)

	r.Header.Set("Content-Type", "application/json")

	return r
}
