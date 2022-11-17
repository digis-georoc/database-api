package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/handler"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/repository"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/secretstore"

	// import swagger docs
	_ "gitlab.gwdg.de/fe/digis/database-api/cmd/docs"
)

func main() {
	// setup database connection
	db := repository.NewPostgresConnector()
	defer db.Close()

	// for local testing this can be set to the local project directory; in containerized setup this remains empty
	workdir := os.Getenv("WORKDIR")
	secStore := secretstore.NewSecretStore(workdir)
	err := secStore.LoadSecretsFromFile("/vault/secrets/database-config.txt")
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

	handler := handler.NewHandler(db, nil)
	echoAPI := api.InitializeAPI(handler, secStore)

	// start api server
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "80"
	}
	log.Fatal(echoAPI.Start(":" + port))
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

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, database), nil
}
