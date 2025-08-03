package lib

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetHTTPData(path string) ([]byte, error) {
	response, err := http.Get(path)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return responseData, nil
}
