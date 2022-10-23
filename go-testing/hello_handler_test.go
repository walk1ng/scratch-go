package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_hellohandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.io/hello", nil)
	w := httptest.NewRecorder()

	hellohandler(w, req)
	b, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err.Error())
	}
	if string(b) != "hello world!!" {
		t.Fatal("expected hello world!!, but got", string(b))
	}
}
