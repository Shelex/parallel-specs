package mock

import (
	"github.com/Shelex/split-specs-v2/internal/appError"
	"github.com/Shelex/split-specs-v2/internal/entities"
	"github.com/google/uuid"
)

func (i *MockStorage) DeleteSpecs(projectID string) error {
	for _, spec := range i.Specs {
		if spec.ProjectID == projectID {
			delete(i.Specs, spec.ID)
		}
	}
	return nil
}

func (i *MockStorage) AddSpecsMaybe(projectID string, names []string) ([]entities.Spec, error) {
	specs := make([]entities.Spec, len(names))
	for index, name := range names {
		spec, err := i.IsSpecAvailable(projectID, name)
		if err != nil {
			if err == appError.SpecNotFound {
				newSpec, err := i.AddSpec(projectID, name)
				if err != nil {
					return nil, err
				}
				spec = newSpec
			} else {
				return nil, err
			}
		}
		specs[index] = spec
	}
	return specs, nil
}

func (i *MockStorage) IsSpecAvailable(projectID string, name string) (entities.Spec, error) {
	for _, spec := range i.Specs {
		if spec.FilePath == name {
			return *spec, nil
		}
	}
	return entities.Spec{}, appError.SpecNotFound
}

func (i *MockStorage) AddSpec(projectID string, name string) (entities.Spec, error) {
	id := uuid.NewString()
	spec := entities.Spec{
		ID:        id,
		FilePath:  name,
		ProjectID: projectID,
	}
	i.Specs[id] = &spec
	return spec, nil
}

func (i *MockStorage) GetSpec(specID string) (entities.Spec, error) {
	spec, ok := i.Specs[specID]

	if !ok {
		return entities.Spec{}, appError.SpecNotFound
	}

	return *spec, nil
}
