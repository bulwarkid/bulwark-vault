//go:build js && wasm

package sdk

import (
	"errors"
	"fmt"
	"syscall/js"
)

const BASE_URL = "http://localhost:8080"

type UnsuccessfulRequestError struct {
	code     int
	response string
}

func (e *UnsuccessfulRequestError) Error() string {
	return fmt.Sprintf("Unsuccessful request (%d): %s", e.code, e.response)
}

func requestError(code int, response string) *UnsuccessfulRequestError {
	return &UnsuccessfulRequestError{code, response}
}

func isReturnCode(err error, code int) bool {
	var requestError *UnsuccessfulRequestError
	return errors.As(err, &requestError) && requestError.code == code
}

func request(method, path, data string) (string, error) {
	request := js.Global().Get("wasmApi").Get("request")
	output, _ := await(request.Invoke(js.ValueOf(method), js.ValueOf(BASE_URL + path), data))
	if output != nil {
		responseText := output.Get("data").String()
		statusCode := output.Get("code").Int()
		if statusCode >= 300 || statusCode < 200 {
			return "", fmt.Errorf("%w", requestError(statusCode, responseText))
		}
		return responseText, nil
	}
	return "", fmt.Errorf("Could not read response")
}

func post(path, dataType, data string) (string, error) {
	return request("POST", path, data)
}

func get(path string) (string, error) {
	return request("GET", path, "")
}
