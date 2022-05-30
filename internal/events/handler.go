package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type Orchestrator struct {
	Channel *gochannel.GoChannel
}

var Handler Orchestrator

func Start() *Orchestrator {
	logger := watermill.NewStdLogger(true, true).With(
		watermill.LogFields{},
	)
	Handler = Orchestrator{
		Channel: gochannel.NewGoChannel(
			gochannel.Config{
				Persistent: false,
			},
			logger,
		),
	}

	return &Handler
}

func eventToMessage(topic TopicName, event interface{}) *message.Message {
	marshalled, _ := json.Marshal(event)
	return &message.Message{
		Payload: marshalled,
		Metadata: map[string]string{
			"topic": topic.String(),
		},
	}
}

const (
	Events = "events"
)

func (e *Orchestrator) Publish(topic TopicName, event interface{}) {
	if err := e.Channel.Publish(Events, eventToMessage(topic, event)); err != nil {
		log.Printf("[events] failed to publish event %s", err)
	}
}

func (e *Orchestrator) Subscribe() (<-chan *message.Message, error) {
	messages, err := e.Channel.Subscribe(context.Background(), Events)
	if err != nil {
		return nil, fmt.Errorf("[events] failed to subscribe: %s", err)
	}

	return messages, nil
}
