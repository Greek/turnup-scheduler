package http

import (
	"io"
	"net/http"
)

func GetHTTPData(path string) ([]byte, error) {
	response, err := http.Get(path)
	if err != nil {
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}
