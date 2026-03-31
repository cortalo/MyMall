package shared

import (
	"MyMall/domain"
	"context"
)

type EventPublisher interface {
	Publish(ctx context.Context, event domain.Event) error
}
