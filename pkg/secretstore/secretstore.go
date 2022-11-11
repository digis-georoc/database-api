package secretstore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// secretStore holds a map of secrets provided by an external source, e.g. Vault
// secretMap is a key/value map containing the secrets
type secretStore struct {
	secretMap map[string]string
}

// SecretStore provides the interface methods to interact with a secretStore
type SecretStore interface {
	// Loads a set of key/value pairs from a file into the secretMap of the secretStore instance
	LoadSecretsFromFile(path string) error

	// Returns a secret by its key or nil if the key is not found
	GetSecret(key string) (string, error)
}

// NewSecretStore creates a new secretStore from a given file path
func NewSecretStore(path string) (*secretStore, error) {
	store := secretStore{}
	err := store.LoadSecretsFromFile(path)
	return &store, err
}

// LoadSecretsFromFile loads the secrets from the given filepath
// Is invoked on NewSecretstore(path)
func (s *secretStore) LoadSecretsFromFile(path string) error {
	secretMap := make(map[string]string)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	secretBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(secretBytes, &secretMap)
	if err != nil {
		return err
	}

	s.secretMap = secretMap
	return nil
}

func (s *secretStore) GetSecret(key string) (string, error) {
	secret, ok := s.secretMap[key]
	if !ok {
		return "", fmt.Errorf("No secret with key '%s'", key)
	}
	return secret, nil
}
