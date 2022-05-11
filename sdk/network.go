package sdk

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const BASE_URL = "http://localhost:8080"

type UnsuccessfulRequestError struct {
	code     int
	response string
}

func (e *UnsuccessfulRequestError) Error() string {
	return fmt.Sprintf("Unsuccessful request (%d): %s", e.code, e.response)
}

func isReturnCode(err error, code int) bool {
	var requestError *UnsuccessfulRequestError
	return errors.As(err, &requestError) && requestError.code == code
}

func get(path string) (string, error) {
	response, err := http.Get(BASE_URL + path)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if response.StatusCode != 200 {
		return "", fmt.Errorf("%w", &UnsuccessfulRequestError{code: response.StatusCode, response: string(bytes)})
	}
	return string(bytes), nil
}

func post(path string, dataType string, data string) (string, error) {
	response, err := http.Post(BASE_URL+path, dataType, strings.NewReader(data))
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
