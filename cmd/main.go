// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/api/handler"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/repository"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/secretstore"
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
	params, err := getConnectionParams(secStore)
	if err != nil {
		log.Fatal(fmt.Errorf("Can not get connection params: %v", err))
	}
	if params.SSHHost != "" && params.SSHPort != 0 {
		err = db.SSHConnect(connString, params)
		if err != nil {
			log.Fatal(fmt.Errorf("Can not connect to database via ssh: %w", err))
		}
	} else {
		err = db.Connect(connString)
		if err != nil {
			log.Fatal(fmt.Errorf("Can not connect to database: %w", err))
		}
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(fmt.Errorf("Can not reach database: %w", err))
	}
	log.Infof("Connected to database: %s/%s", params.DBHost, params.DBName)

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
	params, err := getConnectionParams(secStore)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", params.DBUser, params.DBPassword, params.DBHost, params.DBPort, params.DBName), nil
}

// getConnectionParams retrieves the database connection parameters from vault- and env-vars
// param secStore: the instance of the secretstore.Secretstore to load values provided by vault
func getConnectionParams(secStore secretstore.SecretStore) (*repository.ConnectionParams, error) {
	username, err := secStore.GetSecret("DB_USER")
	if err != nil {
		return nil, fmt.Errorf("Can not load secret DB_USER")
	}
	password, err := secStore.GetSecret("DB_PASSWORD")
	if err != nil {
		return nil, fmt.Errorf("Can not load secret DB_PASSWORD")
	}

	// SSH params are optional
	sshUser, _ := secStore.GetSecret("SSH_USER")
	sshPassword, _ := secStore.GetSecret("SSH_PASSWORD")

	host := os.Getenv("DB_HOST")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("Can not cast env-var DB_PORT to int: %+v", os.Getenv("DB_PORT"))
	}
	database := os.Getenv("DB_NAME")
	sshHost := os.Getenv("SSH_HOST")
	sshPort := 0
	sshPortEnv := os.Getenv("SSH_PORT")
	if sshPortEnv != "" {
		sshPort, err = strconv.Atoi(sshPortEnv)
		if err != nil {
			return nil, fmt.Errorf("Can not cast env-var SSH_PORT to int: '%+v'", os.Getenv("SSH_PORT"))
		}
	}

	return &repository.ConnectionParams{DBHost: host, DBPort: port, DBUser: username, DBPassword: password, DBName: database, SSHHost: sshHost, SSHPort: sshPort, SSHUser: sshUser, SSHPassword: sshPassword}, nil
}
