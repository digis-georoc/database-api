package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gitlab.gwdg.de/fe/digis/database-api/src/api"
	"gitlab.gwdg.de/fe/digis/database-api/src/repository"
)

func main() {
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

	authors, err := db.GetAuthorByName("smith")
	if err != nil {
		log.Errorf("Can not retrieve authors: %v", err)
	}
	log.Infof("Retrieved %d authors: %v", len(authors), authors)

	echoAPI := api.InitializeAPI()
	log.Fatal(echoAPI.Start(":81"))
}
