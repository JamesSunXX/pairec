package web

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestControllerReadRequestBodyWithGzip(t *testing.T) {
	ctrl := Controller{}
	want := []byte(`{"foo":"bar"}`)

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(want); err != nil {
		t.Fatalf("gzip write: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("gzip close: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(buf.Bytes()))
	req.Header.Set("Content-Encoding", "gzip")

	got, err := ctrl.ReadRequestBody(req)
	if err != nil {
		t.Fatalf("ReadRequestBody returned error: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("unexpected body: got %s, want %s", string(got), string(want))
	}
}

func TestControllerResponseWithGzip(t *testing.T) {
	ctrl := Controller{}
	body := []byte(`{"msg":"hello"}`)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rr := httptest.NewRecorder()

	ctrl.Response(rr, req, body)

	res := rr.Result()
	if res.Header.Get("Content-Encoding") != "gzip" {
		t.Fatalf("Content-Encoding not set to gzip, got %q", res.Header.Get("Content-Encoding"))
	}

	gz, err := gzip.NewReader(res.Body)
	if err != nil {
		t.Fatalf("create gzip reader: %v", err)
	}
	defer gz.Close()

	got, err := io.ReadAll(gz)
	if err != nil {
		t.Fatalf("read gzip body: %v", err)
	}
	if !bytes.Equal(got, body) {
		t.Fatalf("unexpected response body: got %s, want %s", string(got), string(body))
	}
}
