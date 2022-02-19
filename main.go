package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/njoysubho/whatson/client"
	"github.com/njoysubho/whatson/dto"
	"github.com/njoysubho/whatson/handler"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	tmdbRequest := &dto.TmdbRequest{}
	err := json.Unmarshal([]byte(request.Body), tmdbRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid request Body"}
	}
	apiKey, _ := client.GetSecret()
	tmdbApi := client.TMDBApiClient{BaseUrl: "https://api.themoviedb.org/3",
		ApiKey: apiKey}
	queryExecutor:=handler.QueryExecutor{Apiclient: tmdbApi}
	response := queryExecutor.QueryTmdb(tmdbRequest)
	return response
}

func main() {
	lambda.Start(HandleRequest)
}
