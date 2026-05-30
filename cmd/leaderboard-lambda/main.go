package main

import (
	"context"
	"log"

	"snake_golang/internal/leaderboard"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	handler, err := leaderboard.NewHandler(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return handler.Handle(ctx, req)
	})
}
