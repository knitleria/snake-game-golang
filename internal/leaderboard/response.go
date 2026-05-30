package leaderboard

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func jsonResponse(status int, body any) (events.APIGatewayV2HTTPResponse, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, fmt.Errorf("marshal response: %w", err)
	}
	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers: map[string]string{
			"content-type": "application/json; charset=utf-8",
		},
		Body: string(jsonBody),
	}, nil
}

func errorResponse(status int, message string) (events.APIGatewayV2HTTPResponse, error) {
	return jsonResponse(status, ErrorResponse{Error: message})
}
