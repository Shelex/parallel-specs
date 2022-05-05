package mock

import (
	"fmt"

	"github.com/Shelex/split-specs-v2/internal/appError"
	"github.com/Shelex/split-specs-v2/internal/entities"
)

func (i *MockStorage) AddApiKey(userID string, key entities.ApiKey) error {
	_, ok := i.ApiKeys[key.ID]
	if ok {
		return fmt.Errorf("api key with id %s already exist", key.ID)
	}

	i.ApiKeys[key.ID] = &key

	return nil
}

func (i *MockStorage) DeleteApiKey(userID string, keyID string) error {
	_, ok := i.ApiKeys[keyID]
	if !ok {
		return appError.ApiKeyNotFound
	}

	delete(i.ApiKeys, keyID)
	return nil
}

func (i *MockStorage) GetApiKeys(userID string) ([]entities.ApiKey, error) {
	var keys []entities.ApiKey

	for _, key := range i.ApiKeys {
		if key.UserID == userID {
			keys = append(keys, *key)
		}
	}

	return keys, nil
}

func (i *MockStorage) IsApiKeyValid(userID string, keyID string) (bool, error) {
	for _, key := range i.ApiKeys {
		if key.ID == keyID {
			return true, nil
		}
	}

	return false, nil
}
