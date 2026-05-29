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
	InteractionEventsTopic   = "interaction.events"
	interactionEventsGroupID = "feed-service-go-interaction-events"
)

func NewInteractionEventsReader() *kgo.Reader {
	return kgo.NewReader(kgo.ReaderConfig{
		Brokers: kafka.Brokers,
		Topic:   InteractionEventsTopic,
		GroupID: interactionEventsGroupID,
		Dialer:  kafka.Dialer,
	})
}

func RunInteractionEventsConsumer(ctx context.Context, interactionSvc service.IInteractionService) error {
	r := NewInteractionEventsReader()
	defer func() { _ = r.Close() }()

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var event events.InteractionEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("kafka interaction.events: invalid json partition=%d offset=%d: %v", m.Partition, m.Offset, err)
			continue
		}

		eventType := strings.TrimSpace(event.EventType)
		actorUserID := strings.TrimSpace(event.ActorUserId)
		targetUserID := strings.TrimSpace(event.TargetUserId)
		postID := strings.TrimSpace(event.PostId)

		if eventType == "" || actorUserID == "" || targetUserID == "" || postID == "" {
			log.Printf("kafka interaction.events: missing required fields partition=%d offset=%d", m.Partition, m.Offset)
			continue
		}

		if err := interactionSvc.ApplyInteractionEvent(ctx, event); err != nil {
			log.Printf("kafka interaction.events: apply failed postID=%s actor=%s: %v", postID, actorUserID, err)
			continue
		}

		log.Printf("kafka interaction.events: ok eventType=%s postID=%s actor=%s partition=%d offset=%d",
			eventType, postID, actorUserID, m.Partition, m.Offset)
	}
}
