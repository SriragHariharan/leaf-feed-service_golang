package repo

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type feedCursor struct {
	CreatedAt time.Time `json:"created_at"`
	FeedID    string    `json:"feed_id"`
}

func encodeFeedCursor(c feedCursor) (string, error) {
	payload, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(payload), nil
}

func decodeFeedCursor(cursor string) (*feedCursor, error) {
	if cursor == "" {
		return nil, nil
	}

	payload, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	var decoded feedCursor
	if err := json.Unmarshal(payload, &decoded); err != nil {
		return nil, err
	}
	return &decoded, nil
}
