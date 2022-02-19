package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
)
type Client interface {
	doGET(uri string) ([]byte, error)
}

type HttpClient struct{}

func (client *HttpClient) doGET(uri string) ([]byte, error) {
	apiResponse, err := http.Get(uri)

	if err != nil {
		return nil, fmt.Errorf("Unable to make call")
	}
	defer apiResponse.Body.Close()

	body, err := ioutil.ReadAll(apiResponse.Body)

	if err != nil {
		panic("Unable to parse")
	}
	return body, nil
}

