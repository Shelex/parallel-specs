package execution

import (
	"github.com/Shelex/split-specs-v2/internal/entities"
	"github.com/Shelex/split-specs-v2/repository"
	"github.com/google/uuid"
)

func SpecsToExecutions(specs []entities.Spec) ([]entities.Execution, error) {
	specExecutions := make([]entities.Execution, len(specs))

	for index, spec := range specs {
		execution, err := specToExecution(spec)
		if err != nil {
			return nil, err
		}
		specExecutions[index] = *execution
	}
	return specExecutions, nil
}

func specToExecution(spec entities.Spec) (*entities.Execution, error) {
	executions, err := repository.DB.GetExecutionHistory(spec.ID, 5)
	if err != nil {
		return nil, err
	}

	execution := entities.Execution{
		ID:       uuid.NewString(),
		SpecID:   spec.ID,
		SpecName: spec.FilePath,
	}

	if len(executions) == 0 {
		return &execution, nil
	}

	var sum uint32

	for _, execution := range executions {
		sum += execution.Duration
	}

	average := int(sum) / len(executions)

	execution.EstimatedDuration = uint32(average)
	return &execution, nil
}
