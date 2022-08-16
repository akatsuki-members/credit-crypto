package subscribers

import (
	"context"

	"github.com/akatsuki-members/credit-crypto/libs/pubsub/messages"
	"github.com/pkg/errors"
)

type EventBusSubscriber interface {
	// Pull pull a group of messages from event bus.
	Pull(ctx context.Context, numberOfMessages uint8) ([]messages.Event, error)
	Stream(ctx context.Context) (<-chan messages.Event, error)
	Subscribe(ctx context.Context, channel string) error
	Acknowledge(ctx context.Context, ID string) error
}

type Subscriber struct {
	eventBus        EventBusSubscriber
	channel         string
	messagesPerPull uint8
}

type Settings struct {
	EventBus        EventBusSubscriber
	MessagesPerPull uint8
}

var errNoChannelName = errors.New("must provide a channel name")

func New(settings Settings) *Subscriber {
	newSubscriber := Subscriber{
		eventBus:        settings.EventBus,
		channel:         "",
		messagesPerPull: settings.MessagesPerPull,
	}

	return &newSubscriber
}

func (s *Subscriber) Subscribe(ctx context.Context, channel string) error {
	if channel == "" {
		return errNoChannelName
	}

	err := s.eventBus.Subscribe(ctx, channel)
	if err != nil {
		return errors.WithMessagef(err, "unexpected error subscribing to channel %q", channel)
	}

	s.channel = channel

	return nil
}

// Pull calls the event bus pull method to get the configured number of messages.
func (s *Subscriber) Pull(ctx context.Context) ([]messages.Event, error) {
	result, err := s.eventBus.Pull(ctx, s.messagesPerPull)
	if err != nil {
		return nil, errors.WithMessagef(err, "unexpected error pulling messages from %q", s.channel)
	}

	return result, nil
}

// Stream calls the event bus stream method to get a stream of events.
func (s *Subscriber) Stream(ctx context.Context) (<-chan messages.Event, error) {
	stream, err := s.eventBus.Stream(ctx)
	if err != nil {
		return nil, errors.WithMessagef(err, "unexpected error streaming messages from %q", s.channel)
	}

	return stream, nil
}

// Acknowledge acknowledge a given message id to avoid processing it again.
func (s *Subscriber) Acknowledge(ctx context.Context, messageID string) error {
	err := s.eventBus.Acknowledge(ctx, messageID)
	if err != nil {
		return errors.WithMessagef(err, "unexpected error acknowledging message: %q", messageID)
	}

	return nil
}
