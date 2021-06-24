package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	jsonStr := `
		{
			"hey": "yo!"
		}
	`
	var mapJson map[string]string
	json.Unmarshal([]byte(jsonStr), &mapJson)
	req := events.APIGatewayProxyRequest{
		QueryStringParameters: mapJson,
	}

	testName := "test1"
	t.Run(testName, func(t *testing.T) {
		ret, err := handler(req)
		if err != nil {
			t.Errorf(testName)
		}
		fmt.Printf("StatusCode=%d, Body=%s\n", ret.StatusCode, ret.Body)
	})
}
