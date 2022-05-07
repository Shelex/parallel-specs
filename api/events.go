package api

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/Shelex/split-specs-v2/internal/events"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gofiber/websocket/v2"
)

func (c *Controller) Listener(conn *websocket.Conn) {
	userId := conn.Query("userId", "")

	if userId == "" {
		return
	}

	projectId := conn.Query("projectId", "")
	sessionId := conn.Query("sessionId", "")

	sessions := make(<-chan *message.Message)
	executions := make(<-chan *message.Message)
	projects := make(<-chan *message.Message)

	if userId != "" {
		channel, err := c.app.Events.Subscribe(events.Project)
		if err != nil {
			log.Printf("failed to subscribe to topic: %s", err)
		}
		projects = channel
	}

	if projectId != "" {
		channel, err := c.app.Events.Subscribe(events.Session)
		if err != nil {
			log.Printf("failed to subscribe to topic: %s", err)
		}
		sessions = channel
	}

	if sessionId != "" {
		channel, err := c.app.Events.Subscribe(events.Execution)
		if err != nil {
			log.Printf("failed to subscribe to topic: %s", err)
		}
		executions = channel
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		for project := range projects {
			var event events.ProjectEvent
			if err := json.Unmarshal(project.Payload, &event); err != nil {
				log.Printf("failed to unmarshall event payload: %s", err)
			}
			if event.UserID == userId {
				event.Topic = events.Project
				if err := conn.WriteJSON(event); err != nil {
					log.Printf("failed to send event: %s", err)
				}
				project.Ack()
			}
		}
	}()

	go func() {
		defer wg.Done()
		for session := range sessions {
			var event events.SessionEvent
			if err := json.Unmarshal(session.Payload, &event); err != nil {
				log.Printf("failed to unmarshall event payload: %s", err)
			}
			if event.ID == sessionId {
				event.Topic = events.Session
				if err := conn.WriteJSON(event); err != nil {
					log.Printf("failed to send event: %s", err)
				}
				session.Ack()
			}
		}
	}()

	go func() {
		defer wg.Done()
		for execution := range executions {
			var event events.ExecutionEvent
			if err := json.Unmarshal(execution.Payload, &event); err != nil {
				log.Printf("failed to unmarshall event payload: %s", err)
			}
			if event.SessionID == sessionId {
				event.Topic = events.Execution
				if err := conn.WriteJSON(event); err != nil {
					log.Printf("failed to send event: %s", err)
				}
				execution.Ack()
			}
		}
	}()

	wg.Wait()
}
