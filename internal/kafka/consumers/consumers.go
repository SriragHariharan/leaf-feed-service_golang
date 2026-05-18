package consumers

import (
	"context"
	"errors"
	"log"

	"github.com/sriraghariharan/feed-service-go/internal/service"
)

func Run(ctx context.Context, userSvc service.IUserService) {
	go func() {
		if err := RunUserEventsConsumer(ctx, userSvc); err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("user.events consumer: %v", err)
		}
	}()
}
