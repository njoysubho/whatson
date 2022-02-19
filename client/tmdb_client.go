package client

import (
	"encoding/json"
	"fmt"
	"sync"
)

type TMDBApiClient struct {
	BaseUrl string
	HttpClient  Client
	ApiKey  string
}

type TMDBTrendResponse struct {
	Page    int      `json:"page"`
	Results []Result `json:"results"`
}

type Result struct {
	BackDropPath     string  `json:"backdrop_path"`
	Id               int     `json:"id"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	Popularity       float64 `json:"popularity"`
}

func (tmdb *TMDBApiClient) DiscoverAll(maxResultCount int) (resp *TMDBTrendResponse, e error) {
	combinedResponse := &TMDBTrendResponse{}
	movies := make([]Result, 0)
	tvs := make([]Result, 0)
	var wg sync.WaitGroup
	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		response, err := tmdb.discoverMovies(tmdb.BaseUrl + "/movie/popular")
		if err != nil {
			fmt.Printf("error %s while fetching movies", err)
		}
		fmt.Printf("Found %d movies", len(response.Results))
		movies = appendUptoMax(response, movies,maxResultCount)
	}(&wg)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		response, err := tmdb.discoverTV(tmdb.BaseUrl + "/tv/popular")
		if err != nil {
			fmt.Printf("error %s while fetching tvs", err)
		}
		fmt.Printf("Found %d tvShows", len(response.Results))
		tvs = appendUptoMax(response, tvs,maxResultCount)
	}(&wg)
	wg.Wait()
	combinedResponse.Results = append(movies, tvs...)
	return combinedResponse, nil
}

func appendUptoMax(response *TMDBTrendResponse, results []Result, max int) []Result {
	if len(response.Results) < max{
		results = append(results, response.Results...)
	} else {
		results = append(results, response.Results[0:max]...)
	}
	return results
}

func (tmdb *TMDBApiClient) discoverMovies(uri string) (resp *TMDBTrendResponse, e error) {
	response, err := tmdb.HttpClient.doGET(uri+"?api_key=" + tmdb.ApiKey)
	err = json.Unmarshal(response, &resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (tmdb *TMDBApiClient) discoverTV(uri string) (resp *TMDBTrendResponse, e error) {
	response, err :=  tmdb.HttpClient.doGET(uri+"?api_key=" + tmdb.ApiKey)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(response, &resp)
	return resp, err
}

