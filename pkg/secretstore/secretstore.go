package secretstore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// secretStore holds a map of secrets provided by an external source, e.g. Vault
// secretMap is a key/value map containing the secrets
// workdir is the current directory relative to which the paths are evaluated in calls to LoadSecretsFromFile(path)
type secretStore struct {
	secretMap map[string]string
	workdir   string
}

// SecretStore provides the interface methods to interact with a secretStore
type SecretStore interface {
	// Loads a set of key/value pairs from a file into the secretMap of the secretStore instance
	LoadSecretsFromFile(path string) error

	// Returns a secret by its key or nil if the key is not found
	GetSecret(key string) (string, error)

	// Returns the full secret map
	GetMap() (map[string]string, error)
}

// NewSecretStore creates a new secretStore from a given file path
func NewSecretStore(workdir string) SecretStore {
	store := secretStore{
		workdir: workdir,
	}
	return &store
}

// LoadSecretsFromFile loads the secrets from the given filepath
func (s *secretStore) LoadSecretsFromFile(path string) error {
	secretMap := make(map[string]string)
	fullPath := s.workdir + path
	f, err := os.Open(fullPath)
	if err != nil {
		return fmt.Errorf("Can not open path %s: %v", fullPath, err)
	}
	defer f.Close()

	secretBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("Can not read file: %v", err)
	}
	err = json.Unmarshal(secretBytes, &secretMap)
	if err != nil {
		return fmt.Errorf("Can not unmarshal data: %v", err)
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

func (s *secretStore) GetMap() (map[string]string, error) {
	if s.secretMap == nil || len(s.secretMap) == 0 {
		return map[string]string{}, fmt.Errorf("No secretMap found or empty: %+v", s.secretMap)
	}
	return s.secretMap, nil
}
