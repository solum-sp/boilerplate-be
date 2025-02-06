package service

import (
	"context"
)

// region:      ======= interface =======
type Publisher interface {
	SendMessage(ctx context.Context, value interface{}) error
}

type Subscriber interface {
	SubscribeToTopic(ctx context.Context) error
	ConsumeMessages(ctx context.Context, msgTypeConf func() ConsumerMessage) (chMsg <-chan ConsumerMessage, chErr <-chan error, chCommitRequest chan<- bool)
}
type ConsumerMessage interface {
	EventName() string
}

// endregion:   ======= interface =======
