package entities

type User struct {
	ID       string `json:"id" swaggerignore:"true"`
	Email    string `json:"email" validate:"required,email,min=6,max=32"`
	Password string `json:"password" validate:"required,min=3,max=32"`
}

type UserProject struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	ProjectID string `json:"projectId"`
}

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Session struct {
	ID         string      `json:"id"`
	ProjectID  string      `json:"projectId"`
	Executions []Execution `json:"executions,omitempty"`
	StartedAt  uint64      `json:"startedAt"`
	FinishedAt uint64      `json:"finishedAt"`
	CreatedAt  uint64      `json:"createdAt"`
}

type Spec struct {
	ID        string `json:"id"`
	FilePath  string `json:"filePath"`
	ProjectID string `json:"projectId"`
}

type Execution struct {
	ID                string `json:"id"`
	SpecID            string `json:"specId"`
	SpecName          string `json:"specName"`
	SessionID         string `json:"sessionId"`
	MachineID         string `json:"machineId"`
	StartedAt         uint64 `json:"startedAt"`
	FinishedAt        uint64 `json:"finishedAt"`
	EstimatedDuration uint32 `json:"estimatedDuration"`
	Duration          uint32 `json:"duration"`
	Status            string `json:"status"`
}

type ApiKey struct {
	ID       string `json:"id"`
	UserID   string `json:"userId"`
	Name     string `json:"name"`
	ExpireAt uint64 `json:"expireAt"`
}

type Pagination struct {
	Limit  int `json:"limit" query:"limit"`
	Offset int `json:"offset" query:"offset"`
}
