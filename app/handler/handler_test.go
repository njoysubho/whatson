package handler

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/njoysubho/whatson/app/client"
	"github.com/njoysubho/whatson/app/dto"
	"testing"
)

var baseUrl = "https://api.themoviedb.org/3"

func TestItShouldGetDefaultQuery(t *testing.T) {
	tmdbRequest := &dto.TmdbRequest{}
	request := events.APIGatewayProxyRequest{Body: "{\"query\": \"default\"}"}
	err := json.Unmarshal([]byte(request.Body), &tmdbRequest)
	if err != nil {
		t.Errorf("Error while unmarshaliing request body")
	}
	var mockHttpClient = client.MockHttpClient{Response: map[client.MockRequest]interface{}{}}
	mockResponse := mockTmdbTrendResponse(5)

	mockHttpClient.StubGet(baseUrl+"/movie/popular?api_key=test", mockResponse)
	mockHttpClient.StubGet(baseUrl+"/tv/popular?api_key=test", mockResponse)

	tmdbClient := client.TMDBApiClient{BaseUrl: baseUrl, HttpClient: &mockHttpClient, ApiKey: "test"}

	queryExecutor := QueryExecutor{Apiclient: tmdbClient}
	response := queryExecutor.QueryTmdb(tmdbRequest)
	if response.StatusCode != 200 {
		t.Errorf("Expected status code 200 Found %d", response.StatusCode)
	}
	if response.Body == "" {
		t.Errorf("Expected response body not empty")
	}
}

func mockTmdbTrendResponse(count int) interface{} {
	result := make([]client.Result, 0)
	for i := 0; i < count; i++ {
		result = append(result, client.Result{
			BackDropPath:     "test-backdrop-path",
			Id:               i,
			OriginalLanguage: "en",
			OriginalTitle:    "test-title",
			Popularity:       0,
		})
	}
	return client.TMDBTrendResponse{
		Page:    1,
		Results: result,
	}
}
