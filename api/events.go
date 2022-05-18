package api

import (
	"encoding/json"
	"log"

	"github.com/Shelex/split-specs-v2/internal/events"
	"github.com/gofiber/websocket/v2"
)

func (c *Controller) Listener(conn *websocket.Conn) {
	userId := conn.Query("userId", "")

	defer conn.Close()

	if userId == "" {
		conn.WriteMessage(1, []byte("please specify user id"))
		return
	}

	projectId := conn.Query("projectId", "")
	sessionId := conn.Query("sessionId", "")

	messages, err := c.app.Events.Subscribe()
	if err != nil {
		log.Printf("failed to subscribe to topic: %s", err)
	}

	for message := range messages {
		topic := message.Metadata.Get("topic")
		switch topic {
		case events.Project.String():
			{
				var parsed events.ProjectEvent
				if err := json.Unmarshal(message.Payload, &parsed); err != nil {
					log.Printf("failed to unmarshall event payload: %s", err)
				}
				if parsed.UserID == userId {
					if err := conn.WriteJSON(parsed); err != nil {
						log.Printf("failed to send event: %s", err)
					}
				}
			}
		case events.Session.String():
			{
				var parsed events.SessionEvent
				if err := json.Unmarshal(message.Payload, &parsed); err != nil {
					log.Printf("failed to unmarshall event payload: %s", err)
				}
				if parsed.ProjectID == projectId {
					if err := conn.WriteJSON(parsed); err != nil {
						log.Printf("failed to send event: %s", err)
					}
				}
			}
		case events.Execution.String():
			{
				var parsed events.ExecutionEvent
				if err := json.Unmarshal(message.Payload, &parsed); err != nil {
					log.Printf("failed to unmarshall event payload: %s", err)
				}
				if parsed.SessionID == sessionId {
					if err := conn.WriteJSON(parsed); err != nil {
						log.Printf("failed to send event: %s", err)
					}
				}
			}
		default:
			log.Printf("unknown event topic: %s", topic)
		}
		message.Ack()
	}
}
