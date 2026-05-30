package leaderboard

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Store interface {
	SubmitScore(ctx context.Context, req SubmitScoreRequest, now time.Time) (SubmitScoreResponse, error)
	FetchLeaderboard(ctx context.Context, mode string, limit int) (LeaderboardResponse, error)
}

type DynamoStore struct {
	client    *dynamodb.Client
	tableName string
	indexName string
}

type scoreItem struct {
	PK            string `dynamodbav:"pk"`
	SK            string `dynamodbav:"sk"`
	GSI1PK        string `dynamodbav:"gsi1pk"`
	GSI1SK        string `dynamodbav:"gsi1sk"`
	PlayerID      string `dynamodbav:"player_id"`
	PlayerName    string `dynamodbav:"player_name"`
	Score         int    `dynamodbav:"score"`
	Mode          string `dynamodbav:"mode"`
	Skin          string `dynamodbav:"skin"`
	ClientVersion string `dynamodbav:"client_version"`
	DurationMS    int64  `dynamodbav:"duration_ms"`
	CreatedAt     string `dynamodbav:"created_at"`
	UpdatedAt     string `dynamodbav:"updated_at"`
}

func NewDynamoStore(client *dynamodb.Client, tableName, indexName string) *DynamoStore {
	return &DynamoStore{client: client, tableName: tableName, indexName: indexName}
}

func (s *DynamoStore) SubmitScore(ctx context.Context, req SubmitScoreRequest, now time.Time) (SubmitScoreResponse, error) {
	updatedAt := now.UTC()
	item := scoreItem{
		PK:            playerPK(req.PlayerID),
		SK:            modeSK(req.Mode),
		GSI1PK:        modeGSIKey(req.Mode),
		GSI1SK:        leaderboardSortKey(req.Score, updatedAt, req.PlayerID),
		PlayerID:      req.PlayerID,
		PlayerName:    req.PlayerName,
		Score:         req.Score,
		Mode:          req.Mode,
		Skin:          req.Skin,
		ClientVersion: req.ClientVersion,
		DurationMS:    req.DurationMS,
		UpdatedAt:     updatedAt.Format(time.RFC3339),
	}

	attrs, err := attributevalue.MarshalMap(item)
	if err != nil {
		return SubmitScoreResponse{}, fmt.Errorf("marshal score item: %w", err)
	}

	updateValues := map[string]types.AttributeValue{
		":gsi1pk":         attrs["gsi1pk"],
		":gsi1sk":         attrs["gsi1sk"],
		":player_id":      attrs["player_id"],
		":player_name":    attrs["player_name"],
		":score":          attrs["score"],
		":mode":           attrs["mode"],
		":skin":           attrs["skin"],
		":client_version": attrs["client_version"],
		":duration_ms":    attrs["duration_ms"],
		":updated_at":     attrs["updated_at"],
	}

	out, err := s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: item.PK},
			"sk": &types.AttributeValueMemberS{Value: item.SK},
		},
		ExpressionAttributeNames: map[string]string{
			"#mode": "mode",
		},
		ExpressionAttributeValues: updateValues,
		ConditionExpression:       aws.String("attribute_not_exists(score) OR score < :score"),
		UpdateExpression: aws.String(
			"SET gsi1pk = :gsi1pk, gsi1sk = :gsi1sk, " +
				"player_id = :player_id, player_name = :player_name, score = :score, " +
				"#mode = :mode, skin = :skin, client_version = :client_version, " +
				"duration_ms = :duration_ms, updated_at = :updated_at, " +
				"created_at = if_not_exists(created_at, :updated_at)",
		),
		ReturnValues: types.ReturnValueAllNew,
	})
	if err != nil {
		var conditional *types.ConditionalCheckFailedException
		if errors.As(err, &conditional) {
			best, getErr := s.getBest(ctx, req.PlayerID, req.Mode)
			if getErr != nil {
				return SubmitScoreResponse{}, getErr
			}
			rank, rankErr := s.rank(ctx, best.Mode, best.GSI1SK)
			if rankErr != nil {
				return SubmitScoreResponse{}, rankErr
			}
			return SubmitScoreResponse{
				Accepted: true,
				Improved: false,
				Rank:     rank,
				Entry:    best.entry(rank),
			}, nil
		}
		return SubmitScoreResponse{}, fmt.Errorf("update score: %w", err)
	}

	var saved scoreItem
	if err := attributevalue.UnmarshalMap(out.Attributes, &saved); err != nil {
		return SubmitScoreResponse{}, fmt.Errorf("unmarshal saved score: %w", err)
	}
	rank, err := s.rank(ctx, saved.Mode, saved.GSI1SK)
	if err != nil {
		return SubmitScoreResponse{}, err
	}
	return SubmitScoreResponse{
		Accepted: true,
		Improved: true,
		Rank:     rank,
		Entry:    saved.entry(rank),
	}, nil
}

func (s *DynamoStore) FetchLeaderboard(ctx context.Context, mode string, limit int) (LeaderboardResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}

	out, err := s.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(s.tableName),
		IndexName:              aws.String(s.indexName),
		KeyConditionExpression: aws.String("gsi1pk = :mode"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":mode": &types.AttributeValueMemberS{Value: modeGSIKey(mode)},
		},
		Limit:            aws.Int32(int32(limit)),
		ScanIndexForward: aws.Bool(true),
	})
	if err != nil {
		return LeaderboardResponse{}, fmt.Errorf("query leaderboard: %w", err)
	}

	var items []scoreItem
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &items); err != nil {
		return LeaderboardResponse{}, fmt.Errorf("unmarshal leaderboard: %w", err)
	}

	entries := make([]Entry, 0, len(items))
	for i, item := range items {
		entries = append(entries, item.entry(i+1))
	}
	return LeaderboardResponse{Mode: mode, Entries: entries}, nil
}

func (s *DynamoStore) getBest(ctx context.Context, playerID, mode string) (scoreItem, error) {
	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: playerPK(playerID)},
			"sk": &types.AttributeValueMemberS{Value: modeSK(mode)},
		},
	})
	if err != nil {
		return scoreItem{}, fmt.Errorf("get best score: %w", err)
	}
	if len(out.Item) == 0 {
		return scoreItem{}, fmt.Errorf("best score not found")
	}
	var item scoreItem
	if err := attributevalue.UnmarshalMap(out.Item, &item); err != nil {
		return scoreItem{}, fmt.Errorf("unmarshal best score: %w", err)
	}
	return item, nil
}

func (s *DynamoStore) rank(ctx context.Context, mode, currentSortKey string) (int, error) {
	count := int32(0)
	input := &dynamodb.QueryInput{
		TableName:              aws.String(s.tableName),
		IndexName:              aws.String(s.indexName),
		KeyConditionExpression: aws.String("gsi1pk = :mode AND gsi1sk < :sort_key"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":mode":     &types.AttributeValueMemberS{Value: modeGSIKey(mode)},
			":sort_key": &types.AttributeValueMemberS{Value: currentSortKey},
		},
		Select: types.SelectCount,
	}

	for {
		out, err := s.client.Query(ctx, input)
		if err != nil {
			return 0, fmt.Errorf("query rank: %w", err)
		}
		count += out.Count
		if len(out.LastEvaluatedKey) == 0 {
			break
		}
		input.ExclusiveStartKey = out.LastEvaluatedKey
	}

	return int(count) + 1, nil
}

func (i scoreItem) entry(rank int) Entry {
	return Entry{
		Rank:       rank,
		PlayerName: i.PlayerName,
		Score:      i.Score,
		Mode:       i.Mode,
		Skin:       i.Skin,
		UpdatedAt:  i.UpdatedAt,
	}
}
