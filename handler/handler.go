package handler

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/njoysubho/whatson/client"
	"github.com/njoysubho/whatson/dto"
)
type QueryExecutor struct{
	Apiclient client.TMDBApiClient
}
func (q *QueryExecutor) QueryTmdb(tmdbRequest *dto.TmdbRequest) events.APIGatewayProxyResponse {
	var response events.APIGatewayProxyResponse
	switch tmdbRequest.Query {
	case "movie":
		fmt.Println("Not implemented")
	case "tv":
		fmt.Println("Not implemented")
	default:
		recomendations, err := q.Apiclient.DiscoverAll(4)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Error Getting response from TmDB"}
		}
		jsonResponse, err := json.Marshal(recomendations.Results)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400}
		}
		response = events.APIGatewayProxyResponse{
			StatusCode:      200,
			Body:            string(jsonResponse),
			IsBase64Encoded: false,
		}

	}
	return response
}
