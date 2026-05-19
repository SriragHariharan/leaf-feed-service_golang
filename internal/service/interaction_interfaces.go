package service

import (
	"context"

	"github.com/sriraghariharan/feed-service-go/internal/kafka/events"
)

type IInteractionService interface {
	ApplyInteractionEvent(ctx context.Context, event events.InteractionEvent) error
}
