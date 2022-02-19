package client

import "encoding/json"

type MockRequest struct{
	Url string
	Method string
}
type MockHttpClient struct{
	Response map[MockRequest]interface{}
}

func (client *MockHttpClient) doGET(url string)([]byte,error){
	getRequest := MockRequest{
		Url:    url,
		Method: "GET",
	}
	response := client.Response[getRequest]
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
    return jsonResponse,nil
}

func (client *MockHttpClient) StubGet(url string, response interface{}){
	stubRequest := MockRequest{
		Url:    url,
		Method: "GET",
	}
	client.Response[stubRequest] = response
}





