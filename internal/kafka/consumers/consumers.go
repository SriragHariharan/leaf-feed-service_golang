package consumers

import (
	"context"
	"errors"
	"log"

	"github.com/sriraghariharan/feed-service-go/internal/service"
)

func Run(ctx context.Context, userSvc service.IUserService, interactionSvc service.IInteractionService) {
	go func() {
		if err := RunUserEventsConsumer(ctx, userSvc); err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("user.events consumer: %v", err)
		}
	}()

	go func() {
		if err := RunInteractionEventsConsumer(ctx, interactionSvc); err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("interaction.events consumer: %v", err)
		}
	}()
}
