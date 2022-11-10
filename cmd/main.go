package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/handler"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/middleware"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/repository"
)

func main() {
	// setup database connection
	db := repository.NewPostgresConnector()
	defer db.Close()

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	database := os.Getenv("DATABASE")

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, database)
	err := db.Connect(connString)
	if err != nil {
		log.Fatal(fmt.Errorf("Can not connect to database: %w", err))
	}
	version, err := db.Ping()
	if err != nil {
		log.Fatal(fmt.Errorf("Can not reach database: %w", err))
	}
	log.Infof("Connected to database: %v", version)

	// keycloak configuration
	config := middleware.KeycloakConfig{
		Host:         os.Getenv("KC_HOST"),
		ClientID:     os.Getenv("KC_CLIENTID"),
		ClientSecret: os.Getenv("KC_CLIENTSECRET"),
		Realm:        os.Getenv("KC_REALM"),
	}

	handler := handler.NewHandler(db, config)
	echoAPI := api.InitializeAPI(handler, config)
	log.Fatal(echoAPI.Start(":8081"))
}
