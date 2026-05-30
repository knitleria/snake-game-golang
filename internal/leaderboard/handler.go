package leaderboard

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const maxRequestBodyBytes = 16 * 1024

type Handler struct {
	store    Store
	maxScore int
}

func NewHandler(ctx context.Context) (*Handler, error) {
	tableName := getenv("LEADERBOARD_TABLE", "snake-leaderboard-dev")
	indexName := getenv("LEADERBOARD_INDEX", "LeaderboardIndex")
	maxScore := getenvInt("MAX_SCORE", DefaultMaxScore)

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	var client *dynamodb.Client
	if endpoint := strings.TrimSpace(os.Getenv("DYNAMODB_ENDPOINT")); endpoint != "" {
		client = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	} else {
		client = dynamodb.NewFromConfig(cfg)
	}

	return &Handler{
		store:    NewDynamoStore(client, tableName, indexName),
		maxScore: maxScore,
	}, nil
}

func NewHandlerWithStore(store Store, maxScore int) *Handler {
	return &Handler{store: store, maxScore: maxScore}
}

func (h *Handler) Handle(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	method := req.RequestContext.HTTP.Method
	path := req.RawPath
	if path == "" {
		path = req.RequestContext.HTTP.Path
	}
	if stage := req.RequestContext.Stage; stage != "" && stage != "$default" {
		if trimmed := strings.TrimPrefix(path, "/"+stage); trimmed != path {
			if trimmed == "" {
				trimmed = "/"
			}
			path = trimmed
		}
	}

	switch {
	case method == "GET" && path == "/healthz":
		return jsonResponse(200, map[string]string{"status": "ok"})
	case method == "POST" && path == "/api/v1/scores":
		return h.handleSubmitScore(ctx, req)
	case method == "GET" && path == "/api/v1/leaderboard":
		return h.handleLeaderboard(ctx, req)
	case path == "/healthz" || path == "/api/v1/scores" || path == "/api/v1/leaderboard":
		return errorResponse(405, "method not allowed")
	default:
		return errorResponse(404, "not found")
	}
}

func requestBody(req events.APIGatewayV2HTTPRequest) ([]byte, error) {
	if !req.IsBase64Encoded {
		return []byte(req.Body), nil
	}
	return base64.StdEncoding.DecodeString(req.Body)
}

func (h *Handler) handleSubmitScore(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	body, err := requestBody(req)
	if err != nil {
		return errorResponse(400, "bad request body")
	}
	if len(body) > maxRequestBodyBytes {
		return errorResponse(413, "request body is too large")
	}

	var submit SubmitScoreRequest
	if err := json.Unmarshal(body, &submit); err != nil {
		return errorResponse(400, "bad json")
	}

	submit, err = ValidateSubmitRequest(submit, h.maxScore)
	if err != nil {
		return errorResponse(400, err.Error())
	}

	resp, err := h.store.SubmitScore(ctx, submit, time.Now())
	if err != nil {
		return errorResponse(500, "failed to save score")
	}
	return jsonResponse(200, resp)
}

func (h *Handler) handleLeaderboard(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	mode := strings.TrimSpace(req.QueryStringParameters["mode"])
	if !validMode(mode) {
		return errorResponse(400, ErrBadMode.Error())
	}

	limit := 20
	if raw := strings.TrimSpace(req.QueryStringParameters["limit"]); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return errorResponse(400, "bad limit")
		}
		limit = parsed
	}
	if limit < 1 || limit > 50 {
		return errorResponse(400, "limit is out of range")
	}

	resp, err := h.store.FetchLeaderboard(ctx, mode, limit)
	if err != nil {
		return errorResponse(500, "failed to fetch leaderboard")
	}
	return jsonResponse(200, resp)
}

func getenv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getenvInt(key string, fallback int) int {
	value := getenv(key, "")
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
