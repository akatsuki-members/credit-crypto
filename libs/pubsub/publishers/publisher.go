package publishers

import (
	"context"
	"log"

	"github.com/akatsuki-members/credit-crypto/libs/pubsub/messages"
	"github.com/pkg/errors"
)

// EventBusPublisher defines event bus publishing behavior.
type EventBusPublisher interface {
	// Publish push given event into an event bus in the given message channel.
	Publish(ctx context.Context, messageChannel string, message interface{}) error
}

// EventMessage defines message event and direction.
type EventMessage struct {
	// ChannelName the name of the topic or queue.
	ChannelName string
	// Event event to publish into the event bus.
	Event messages.Event
}

// Publisher define publishing data and logic.
type Publisher struct {
	eventBus EventBusPublisher
}

const (
	publishingErrorMessage = "could not publish event"
)

// New instaces a new publisher.
func New(eventBus EventBusPublisher) *Publisher {
	newEventBus := Publisher{
		eventBus: eventBus,
	}

	return &newEventBus
}

// Publish push given event into the given channel.
func (p *Publisher) Publish(ctx context.Context, event EventMessage) error {
	err := p.eventBus.Publish(ctx, event.ChannelName, event.Event)
	if err != nil {
		log.Println(
			"error", "something went wrong pushing event",
			"content", event,
			"method", "publishers.Publisher.Publish",
		)

		return errors.Wrap(err, publishingErrorMessage)
	}

	return nil
}
