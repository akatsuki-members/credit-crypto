package publishers_test

import (
	"context"
	"testing"

	"github.com/akatsuki-members/credit-crypto/libs/pubsub/messages"
	"github.com/akatsuki-members/credit-crypto/libs/pubsub/publishers"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestPublishSuccess(t *testing.T) {
	t.Parallel()

	// Given
	ctx := context.TODO()
	expectedMessageChannel := "orders-topic"
	expectedEvent := messages.Event{
		Header: messages.Header{
			ID:          "123-456-789",
			Domain:      "loans",
			EventType:   "orders",
			Version:     "0.1.0",
			Application: "core-app",
		},
		Data: []byte(`{"value_one": "one", "value_two": "two"}`),
	}
	event := eventMessageFixture()
	eventBus := new(eventBusMock)
	publisher := publishers.New(eventBus)
	// When
	err := publisher.Publish(ctx, event)
	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedEvent, eventBus.message)
	assert.Equal(t, expectedMessageChannel, eventBus.messageChannel)
}

func TestPublishFailed(t *testing.T) {
	t.Parallel()

	// Given
	ctx := context.TODO()
	expectedError := errors.New("could not publish event: : error")
	expectedMessageChannel := "orders-topic"
	expectedEvent := messages.Event{
		Header: messages.Header{
			ID:          "123-456-789",
			Domain:      "loans",
			EventType:   "orders",
			Version:     "0.1.0",
			Application: "core-app",
		},
		Data: []byte(`{"value_one": "one", "value_two": "two"}`),
	}
	event := eventMessageFixture()
	anError := errors.New("error")
	eventBus := new(eventBusMock).withError(anError)
	publisher := publishers.New(eventBus)
	// When
	err := publisher.Publish(ctx, event)
	// Then
	if !assert.NotNil(t, err) {
		t.Fatalf("expected error but it was nil")
	}

	assert.Equal(t, expectedError.Error(), err.Error())
	assert.Equal(t, expectedEvent, eventBus.message)
	assert.Equal(t, expectedMessageChannel, eventBus.messageChannel)
}

type eventBusMock struct {
	err            error
	messageChannel string
	message        interface{}
}

func (e *eventBusMock) withError(err error) *eventBusMock {
	e.err = err

	return e
}

func (e *eventBusMock) Publish(_ context.Context, channel string, message interface{}) error {
	e.messageChannel = channel
	e.message = message

	return errors.Wrap(e.err, "")
}

func eventMessageFixture() publishers.EventMessage {
	header := messages.Header{
		ID:          "123-456-789",
		Domain:      "loans",
		EventType:   "orders",
		Version:     "0.1.0",
		Application: "core-app",
	}
	event := publishers.EventMessage{
		Event: messages.Event{
			Header: header,
			Data:   []byte(`{"value_one": "one", "value_two": "two"}`),
		},
		ChannelName: "orders-topic",
	}

	return event
}
