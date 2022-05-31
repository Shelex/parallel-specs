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

type BasicEvent struct {
	Topic TopicName `json:"topic"`
	Kind  Name      `json:"kind"`
	ID    string    `json:"id"`
}

type ProjectEvent struct {
	Event  BasicEvent `json:"event"`
	Name   string     `json:"name"`
	UserID string     `json:"userId"`
}

type SessionEvent struct {
	Event     BasicEvent `json:"event"`
	Time      uint64     `json:"time"`
	ProjectID string     `json:"projectId"`
}

type ExecutionEvent struct {
	Event     BasicEvent `json:"event"`
	Time      uint64     `json:"time"`
	SessionID string     `json:"sessionId"`
	MachineID string     `json:"machineId"`
}
