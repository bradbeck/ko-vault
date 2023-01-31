package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/Ralph", nil)
	w := httptest.NewRecorder()
	HelloServer(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	expected := "Hello, Ralph!"
	got := string(data)
	if got != expected {
		t.Errorf("expected %v got %v", expected, got)
	}
}
