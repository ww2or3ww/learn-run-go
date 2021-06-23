package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// AWS Lambda エンドポイント
func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	mapData := map[string]string{}
	if request.QueryStringParameters != nil {
		fmt.Printf("query = %s\n", request.QueryStringParameters)
		mapData = request.QueryStringParameters
	}
	mapData["hello"] = "world!"
	jsonIndent, _ := json.MarshalIndent(mapData, "", "   ")

	return events.APIGatewayProxyResponse{
		Body:       string(jsonIndent),
		StatusCode: 200,
	}, nil
}

// アプリケーションエンドポイント
func main() {
	fmt.Println("=== start main ===")
	lambda.Start(handler)
}
