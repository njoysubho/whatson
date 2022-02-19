package client

import (
	"testing"
)

var baseUrl = "https://api.themoviedb.org/3"


func TestCallToTMDBSuccessfull(t *testing.T) {
	var mockHttpClient = MockHttpClient{Response: map[MockRequest]interface{}{}}
	mockResponse := mockTmdbTrendResponse(5)

	mockHttpClient.StubGet(baseUrl+"/movie/popular?api_key=test", mockResponse)
	mockHttpClient.StubGet(baseUrl+"/tv/popular?api_key=test", mockResponse)

	tmdbClient := TMDBApiClient{BaseUrl: baseUrl, HttpClient: &mockHttpClient, ApiKey: "test"}
	trends, err := tmdbClient.DiscoverAll(4)
	if err != nil {
		t.Errorf("Error calling API")
	}
	if len(trends.Results) != 8 {
		t.Errorf("didn't fetch all  trends ")
	}
}

func TestIgnoreFailedCall(t *testing.T) {
	var mockHttpClient = MockHttpClient{Response: map[MockRequest]interface{}{}}
	mockMovieResponse := mockTmdbTrendResponse(5)
	mockTVResponse := mockTmdbTrendResponse(0)

	mockHttpClient.StubGet(baseUrl+"/movie/popular?api_key=test", mockMovieResponse)
	mockHttpClient.StubGet(baseUrl+"/tv/popular?api_key=test", mockTVResponse)

	tmdbClient := TMDBApiClient{BaseUrl: baseUrl, HttpClient: &mockHttpClient, ApiKey: "test"}
	trends, err := tmdbClient.DiscoverAll(4)
	if err != nil {
		t.Errorf("Error calling API")
	}
	if len(trends.Results) != 4 {
		t.Errorf("didn't fetch all  trends ")
	}
}

func mockTmdbTrendResponse(count int) interface{} {
	result := make([]Result, 0)
	for i := 0; i < count; i++ {
		result = append(result, Result{
			BackDropPath:     "test-backdrop-path",
			Id:               i,
			OriginalLanguage: "en",
			OriginalTitle:    "test-title",
			Popularity:       0,
		})
	}
	return TMDBTrendResponse{
		Page:    1,
		Results: result,
	}
}
