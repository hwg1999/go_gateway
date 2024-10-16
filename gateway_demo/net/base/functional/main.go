package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

type HandleFunc func(http.ResponseWriter, *http.Request)

func (f HandleFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

func main() {
	hf := HandleFunc(HelloHandler)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", bytes.NewBuffer([]byte("test")))
	hf.ServeHTTP(resp, req)

	bts, _ := io.ReadAll(resp.Body)
	fmt.Println(string(bts))
}

func HelloHandler(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello world"))
}
