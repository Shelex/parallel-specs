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
				Persistent: true,
			},
			logger,
		),
	}

	return &Handler
}

func eventToMessage(event interface{}) *message.Message {
	marshalled, _ := json.Marshal(event)

	return &message.Message{
		Payload: marshalled,
	}
}

func (e *Orchestrator) Publish(topic TopicName, event interface{}) {
	if err := e.Channel.Publish(topic.String(), eventToMessage(event)); err != nil {
		log.Printf("[events] failed to publish event topic %s, %s", topic.String(), err)
	}
}

func (e *Orchestrator) Subscribe(topic TopicName) (<-chan *message.Message, error) {
	messages, err := e.Channel.Subscribe(context.Background(), topic.String())
	if err != nil {
		return nil, fmt.Errorf("[events] failed to subscribe to topic %s", topic.String())
	}

	return messages, nil
}
