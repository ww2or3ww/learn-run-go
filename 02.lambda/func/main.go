package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	mapData := map[string]string{}
	mapData["hello"] = "world!"
	jsonIndent, _ := json.MarshalIndent(mapData, "", "   ")
	return events.APIGatewayProxyResponse{
		Body:       string(jsonIndent),
		StatusCode: 200,
	}, nil
}

func main() {
	fmt.Println("=== start main ===")
	lambda.Start(handler)
}
