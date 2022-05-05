package events

type TopicName string

const (
	Project   TopicName = "project"
	Session   TopicName = "session"
	Execution TopicName = "execution"
)

func (t TopicName) String() string {
	return string(t)
}

type Name string

const (
	Created  Name = "created"
	Started  Name = "started"
	Finished Name = "finished"
	Deleted  Name = "deleted"
)

type ProjectEvent struct {
	Topic  TopicName `json:"topic"`
	Kind   Name      `json:"kind"`
	ID     string    `json:"id"`
	Name   string    `json:"name"`
	UserID string    `json:"userId"`
}

type SessionEvent struct {
	Topic     TopicName `json:"topic"`
	Kind      Name      `json:"kind"`
	ID        string    `json:"id"`
	Time      uint64    `json:"time"`
	ProjectID string    `json:"projectId"`
}

type ExecutionEvent struct {
	Topic     TopicName `json:"topic"`
	Kind      Name      `json:"kind"`
	ID        string    `json:"id"`
	Time      uint64    `json:"time"`
	SessionID string    `json:"sessionId"`
}
