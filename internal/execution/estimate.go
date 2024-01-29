package execution

import (
	"fmt"

	"github.com/Shelex/parallel-specs/internal/entities"
	"github.com/Shelex/parallel-specs/internal/errors"
	"github.com/Shelex/parallel-specs/repository"
)

func Next(sessionID string, machineID string, previousStatus string) (string, error) {
	if err := repository.DB.EndExecution(sessionID, machineID, previousStatus); err != nil {
		if err == errors.SpecNotFound {
			return "", errors.SessionNotFound
		}
	}

	allExecutions, err := repository.DB.GetExecutions(sessionID)
	if err != nil {
		return "", err
	}

	executions := getSpecsToRun(allExecutions)

	if len(executions) == len(allExecutions) {
		if err := repository.DB.StartSession(sessionID); err != nil {
			return "", err
		}
	}

	if len(executions) == 0 {
		if err := repository.DB.EndSession(sessionID); err != nil {
			return "", errors.SessionFinished
		}
		return "", errors.SessionFinished
	}

	next := CalculateNext(executions)

	if err := repository.DB.StartExecution(sessionID, machineID, next.SpecID); err != nil {
		return "", fmt.Errorf("failed to start spec: %s", err)
	}

	return next.SpecName, nil
}

func CalculateNext(executions []entities.Execution) entities.Execution {
	newSpec := getNewSpec(executions)
	if newSpec.ID != "" {
		return newSpec
	}

	next := getLongestExecution(executions)
	return next
}

func getSpecsToRun(executions []entities.Execution) []entities.Execution {
	var filtered []entities.Execution
	for _, execution := range executions {
		if execution.StartedAt == 0 {
			filtered = append(filtered, execution)
		}
	}
	return filtered
}

func getLongestExecution(executions []entities.Execution) entities.Execution {
	longestExecution := entities.Execution{}

	for _, execution := range executions {
		if execution.EstimatedDuration > longestExecution.EstimatedDuration {
			longestExecution = execution
		}
	}

	return longestExecution
}

func getNewSpec(executions []entities.Execution) entities.Execution {
	for _, spec := range executions {
		if spec.EstimatedDuration == 0 {
			return spec
		}
	}
	return entities.Execution{}
}
