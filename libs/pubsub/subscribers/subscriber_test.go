package subscribers_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/akatsuki-members/credit-crypto/libs/pubsub/messages"
	"github.com/akatsuki-members/credit-crypto/libs/pubsub/subscribers"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSubscribeAndPull(t *testing.T) {
	t.Parallel()

	// Given
	expectedMessageChannel := "orders-topic"
	expectedEvent := messages.Event{
		Header: messages.Header{
			ID:          "123-456-789",
			Domain:      "loans",
			EventType:   "orders",
			Version:     "0.1.0",
			Application: "core-app",
			MessageID:   "1",
		},
		Data: []byte(`{"value_one": "one", "value_two": "two"}`),
	}
	events := []messages.Event{eventMessageFixture()}
	ctx := context.TODO()
	pullArgs := subscriberData{
		t:               t,
		messagesPerPull: 1,
		eventBus:        new(eventBusMock).withEvents(events),
		channel:         "orders-topic",
	}
	// When
	got, err := subscribeAndPull(ctx, pullArgs)
	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedMessageChannel, pullArgs.eventBus.messageChannel)
	assert.Equal(t, expectedEvent, got[0])
}

func TestSubscribeAndStream(t *testing.T) {
	t.Parallel()

	// Given
	expectedMessageChannel := "orders-topic"
	expectedEvents := eventMessagesFixture()
	events := eventMessagesFixture()
	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)

	defer cancel()

	pullArgs := subscriberData{
		t:                t,
		eventBus:         new(eventBusMock).withEvents(events).withEventStream(),
		channel:          "orders-topic",
		expectedMessages: 2,
	}
	// When
	got, err := subscribeAndStream(ctx, pullArgs)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedMessageChannel, pullArgs.eventBus.messageChannel)
	assert.Equal(t, expectedEvents, got)
}

func TestSubscribeAndPullAndAcknowledge(t *testing.T) {
	t.Parallel()

	// Given
	events := []messages.Event{eventMessageFixture()}
	ctx := context.TODO()
	pullArgs := subscriberData{
		t:               t,
		messagesPerPull: 1,
		eventBus:        new(eventBusMock).withEvents(events),
		channel:         "orders-topic",
	}
	// When
	acknowledgeErrors := pullAndAcknowledge(ctx, pullArgs)
	// Then
	assert.Empty(t, acknowledgeErrors)
}

func TestSubscribeAndError(t *testing.T) {
	t.Parallel()

	// Given
	expectedError := errors.New(`unexpected error subscribing to channel "orders-topic": error`)
	errToReturn := errors.New("error")
	ctx := context.TODO()
	pullArgs := subscriberData{
		t:               t,
		messagesPerPull: 1,
		eventBus:        new(eventBusMock).withSubscribeError(errToReturn),
		channel:         "orders-topic",
	}
	// When
	_, err := subscribeAndPull(ctx, pullArgs)
	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError.Error(), err.Error())
}

func TestPullAndError(t *testing.T) {
	t.Parallel()

	// Given
	expectedError := errors.New(`unexpected error pulling messages from "orders-topic": error`)
	errToReturn := errors.New("error")
	ctx := context.TODO()
	pullArgs := subscriberData{
		t:               t,
		messagesPerPull: 1,
		eventBus:        new(eventBusMock).withPullError(errToReturn),
		channel:         "orders-topic",
	}
	// When
	_, err := subscribeAndPull(ctx, pullArgs)
	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError.Error(), err.Error())
}

func pullAndAcknowledge(ctx context.Context, args subscriberData) []error {
	args.t.Helper()

	settings := subscribers.Settings{
		EventBus:        args.eventBus,
		MessagesPerPull: args.messagesPerPull,
	}
	subscriber := subscribers.New(settings)

	err := subscriber.Subscribe(ctx, args.channel)
	if err != nil {
		args.t.Fatalf("unexpected error: %s", err)
	}

	events, err := subscriber.Pull(ctx)
	if err != nil {
		args.t.Fatalf("unexpected error: %s", err)
	}

	if len(events) == 0 {
		args.t.Fatal("events were expected but got zero")
	}

	var acknowledgeErrors []error

	for idx := range events {
		err := subscriber.Acknowledge(ctx, events[idx].Header.MessageID)
		if err != nil {
			acknowledgeErrors = append(acknowledgeErrors, err)
		}
	}

	return acknowledgeErrors
}

func subscribeAndPull(ctx context.Context, args subscriberData) ([]messages.Event, error) {
	args.t.Helper()

	subscriber, err := subscribe(ctx, args)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	events, err := subscriber.Pull(ctx)
	if args.eventBus.errPull != nil && err == nil {
		args.t.Fatalf("unexpected error: %s", err)
	}

	if args.eventBus.errPull != nil && err != nil {
		return nil, errors.New(err.Error())
	}

	return events, nil
}

func subscribeAndStream(ctx context.Context, args subscriberData) ([]messages.Event, error) {
	args.t.Helper()

	subscriber, err := subscribe(ctx, args)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	stream, err := subscriber.Stream(ctx)
	if args.eventBus.errStream != nil && err == nil {
		args.t.Fatalf("unexpected error: %s", err)
	}

	if args.eventBus.errStream != nil && err != nil {
		return nil, errors.New(err.Error())
	}

	var result []messages.Event

	for i := uint8(0); i < args.expectedMessages; i++ {
		select {
		case <-ctx.Done():
			return result, nil
		case newEvent, ok := <-stream:
			if !ok {
				break
			}

			result = append(result, newEvent)
		}
	}

	return result, nil
}

func subscribe(ctx context.Context, args subscriberData) (*subscribers.Subscriber, error) {
	args.t.Helper()

	settings := subscribers.Settings{
		EventBus:        args.eventBus,
		MessagesPerPull: args.messagesPerPull,
	}
	subscriber := subscribers.New(settings)

	err := subscriber.Subscribe(ctx, args.channel)
	if args.eventBus.errSubscribe != nil && err == nil {
		args.t.Fatalf("unexpected error: %s", err)
	}

	if args.eventBus.errSubscribe != nil && err != nil {
		return nil, errors.New(err.Error())
	}

	return subscriber, nil
}

var errStreamNotInitialized = errors.New("a stream must initialized first")

type subscriberData struct {
	t                *testing.T
	messagesPerPull  uint8
	expectedMessages uint8
	eventBus         *eventBusMock
	channel          string
}

func eventMessageFixture() messages.Event {
	header := messages.Header{
		ID:          "123-456-789",
		Domain:      "loans",
		EventType:   "orders",
		Version:     "0.1.0",
		Application: "core-app",
		MessageID:   "1",
	}
	data := []byte(`{"value_one": "one", "value_two": "two"}`)
	event := messages.Event{
		Header: header,
		Data:   data,
	}

	return event
}

func eventMessagesFixture() []messages.Event {
	events := []messages.Event{
		{
			Header: messages.Header{
				ID:          "123-456-789",
				Domain:      "loans",
				EventType:   "orders",
				Version:     "0.1.0",
				Application: "core-app",
				MessageID:   "1",
			},
			Data: []byte(`{"value_one": "one", "value_two": "two"}`),
		},
		{
			Header: messages.Header{
				ID:          "123-456-790",
				Domain:      "loans",
				EventType:   "orders",
				Version:     "0.1.0",
				Application: "core-app",
				MessageID:   "2",
			},
			Data: []byte(`{"value_one": "three", "value_two": "four"}`),
		},
	}

	return events
}

type eventBusMock struct {
	errPull, errStream, errSubscribe error
	messageChannel                   string
	events                           []messages.Event
	eventStream                      chan messages.Event
}

func (e *eventBusMock) withEventStream() *eventBusMock {
	e.eventStream = make(chan messages.Event)

	return e
}

func (e *eventBusMock) withSubscribeError(err error) *eventBusMock {
	e.errSubscribe = err

	return e
}

func (e *eventBusMock) withPullError(err error) *eventBusMock {
	e.errPull = err

	return e
}

func (e *eventBusMock) withEvents(events []messages.Event) *eventBusMock {
	e.events = events

	return e
}

func (e *eventBusMock) Subscribe(_ context.Context, channel string) error {
	if e.errSubscribe != nil {
		return fmt.Errorf("%w", e.errSubscribe)
	}

	e.messageChannel = channel

	return nil
}

func (e *eventBusMock) Pull(_ context.Context, _ uint8) ([]messages.Event, error) {
	if e.errPull != nil {
		return nil, fmt.Errorf("%w", e.errPull)
	}

	return e.events, nil
}

func (e *eventBusMock) Stream(ctx context.Context) (<-chan messages.Event, error) {
	if e.eventStream == nil {
		return nil, errStreamNotInitialized
	}

	go func() {
		defer close(e.eventStream)

		for _, event := range e.events {
			select {
			case <-ctx.Done():
				return
			case e.eventStream <- event:
			}
		}
	}()

	return e.eventStream, nil
}

func (e *eventBusMock) Acknowledge(ctx context.Context, id string) error {
	return nil
}
