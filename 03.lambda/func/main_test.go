package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandlerSuccess(t *testing.T) {

	req := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"hey": "yo!"},
	}

	testName := "test"
	t.Run(testName, func(t *testing.T) {
		ret, err := handler(req)
		if err != nil {
			t.Errorf(testName)
		}
		fmt.Printf("StatusCode=%d, Body=%s\n", ret.StatusCode, ret.Body)
	})
}
