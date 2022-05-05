package mock

import (
	"github.com/Shelex/split-specs-v2/internal/entities"
	"github.com/Shelex/split-specs-v2/repository"
)

func (i *MockStorage) AddExecutions(sessionID string, executions []entities.Execution) error {
	for _, execution := range executions {
		execution.SessionID = sessionID
		i.Executions[execution.ID] = execution
	}
	return nil
}

func (i *MockStorage) GetExecutions(sessionID string) ([]entities.Execution, error) {
	var executions []entities.Execution
	for _, spec := range i.Executions {
		if spec.SessionID == sessionID {
			executions = append(executions, spec)
		}
	}
	return executions, nil
}

func (i *MockStorage) GetExecutionHistory(specID string, limit int) ([]entities.Execution, error) {
	var executions []entities.Execution
	for _, execution := range i.Executions {
		if execution.SpecID == specID {
			executions = append(executions, execution)
			if len(executions) > limit {
				break
			}
		}
	}
	return executions, nil
}

func (i *MockStorage) StartExecution(sessionID string, machineID string, specID string) error {
	for _, exec := range i.Executions {
		if exec.SpecID == specID && exec.SessionID == sessionID {
			exec.StartedAt = repository.GetTimestamp()
			exec.MachineID = machineID
			i.Executions[exec.ID] = exec
		}
	}
	return nil
}

func (i *MockStorage) EndExecution(sessionID string, machineID string, status string) error {
	for _, exec := range i.Executions {
		if exec.StartedAt > 0 && exec.SessionID == sessionID && exec.MachineID == machineID {
			exec.FinishedAt = repository.GetTimestamp()
			exec.Duration = uint32(exec.FinishedAt - exec.StartedAt)
			exec.Status = status
			i.Executions[exec.ID] = exec
		}
	}
	return nil
}
