package consumers

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	kgo "github.com/segmentio/kafka-go"
	"github.com/sriraghariharan/feed-service-go/internal/kafka"
	"github.com/sriraghariharan/feed-service-go/internal/kafka/events"
	"github.com/sriraghariharan/feed-service-go/internal/service"
)

const (
	UserEventsTopic   = "user.events"
	userEventsGroupID = "feed-service-go-user-events"
)

func NewUserEventsReader() *kgo.Reader {
	return kgo.NewReader(kgo.ReaderConfig{
		Brokers: kafka.Brokers,
		Topic:   UserEventsTopic,
		GroupID: userEventsGroupID,
	})
}

func RunUserEventsConsumer(ctx context.Context, userSvc service.IUserService) error {
	r := NewUserEventsReader()
	defer func() { _ = r.Close() }()

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var event events.UserEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("kafka user.events: invalid json partition=%d offset=%d: %v", m.Partition, m.Offset, err)
			continue
		}

		userID := strings.TrimSpace(event.UserID)
		username := strings.TrimSpace(event.Username)
		if userID == "" || username == "" {
			log.Printf("kafka user.events: missing userID or username partition=%d offset=%d", m.Partition, m.Offset)
			continue
		}

		if err := userSvc.SyncUserFromEvent(ctx, userID, username, event.ProfilePicture); err != nil {
			log.Printf("kafka user.events: sync failed userID=%s: %v", userID, err)
			continue
		}

		log.Printf("kafka user.events: ok userID=%s partition=%d offset=%d", userID, m.Partition, m.Offset)
	}
}
