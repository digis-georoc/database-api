package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/handler"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/repository"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/secretstore"
)

func main() {
	// setup database connection
	db := repository.NewPostgresConnector()
	defer db.Close()

	secStore, err := secretstore.NewSecretStore("/vault/secrets/database-config.txt")
	if err != nil {
		log.Fatal(fmt.Errorf("Error loading secrets: %v", err))
	}
	connString, err := buildConnectionString(secStore)
	if err != nil {
		log.Fatal(fmt.Errorf("Can not build connection string: %v", err))
	}
	err = db.Connect(connString)
	if err != nil {
		log.Fatal(fmt.Errorf("Can not connect to database: %w", err))
	}
	version, err := db.Ping()
	if err != nil {
		log.Fatal(fmt.Errorf("Can not reach database: %w", err))
	}
	log.Infof("Connected to database: %v", version)

	// keycloak configuration
	config, err := buildKeycloakConfig(secStore)
	if err != nil {
		log.Fatal(fmt.Errorf("Can not build keycloak config: %v", err))
	}

	handler := handler.NewHandler(db, *config)
	echoAPI := api.InitializeAPI(handler, *config)
	log.Fatal(echoAPI.Start(":8081"))
}

// buildConnectionString builds the database connection string from vault- and env-vars
// param secStore: the instance of the secretstore.Secretstore to load values provided by vault
func buildConnectionString(secStore secretstore.SecretStore) (string, error) {
	username, err := secStore.GetSecret("DB_USER")
	if err != nil {
		return "", fmt.Errorf("Can not load secret DB_USER")
	}
	password, err := secStore.GetSecret("DB_PASSWORD")
	if err != nil {
		return "", fmt.Errorf("Can not load secret DB_PASSWORD")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	database := os.Getenv("DATABASE")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, database), nil
}

// buildKeycloakConfig builds the configuration for keycloak from vault- and env-vars
// param secStore: the instance of the secretstore.Secretstore to load values provided by vault
func buildKeycloakConfig(secStore secretstore.SecretStore) (*middleware.KeycloakConfig, error) {
	clientID, err := secStore.GetSecret("KC_CLIENTID")
	if err != nil {
		return nil, fmt.Errorf("Can not load secret KC_CLIENTID")
	}
	clientSecret, err := secStore.GetSecret("KC_CLIENTSECRET")
	if err != nil {
		return nil, fmt.Errorf("Can not load secret KC_CLIENTSECRET")
	}
	host := os.Getenv("KC_HOST")
	if host == "" {
		return nil, fmt.Errorf("Can not load env var KC_HOST")
	}
	realm := os.Getenv("KC_REALM")
	if realm == "" {
		return nil, fmt.Errorf("Can not load env var KC_REALM")
	}
	config := middleware.KeycloakConfig{
		Host:         host,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Realm:        realm,
	}
	return &config, nil
}
