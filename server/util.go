package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func returnCode(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf("%d - %s", code, http.StatusText(code))))
}

type httpHandler func(w http.ResponseWriter, r *http.Request)

func requestHandler(getHandler httpHandler, postHandler httpHandler) httpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s - %s: %s\n", r.RemoteAddr, r.Method, r.URL.Path)
		w.Header().Add("Access-Control-Allow-Origin", "*")
		switch r.Method {
		case "GET":
			if getHandler == nil {
				returnCode(w, 405)
			} else {
				getHandler(w, r)
			}
		case "POST":
			if postHandler == nil {
				returnCode(w, 405)
			} else {
				postHandler(w, r)
			}
		default:
			returnCode(w, 405)
		}
	}
}

func readLimit(r io.Reader, limit int64) ([]byte, error) {
	output, err := ioutil.ReadAll(io.LimitReader(r, limit))
	if err != nil {
		return nil, err
	}
	return output, nil
}
