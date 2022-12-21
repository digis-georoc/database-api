package repository

import (
	"context"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PostgresConnector interface exposes methods to connect to and interact with a postgreSQL instance
type PostgresConnector interface {
	// Connect tries to connect to a postgreSQL database with the credentials provided in the connString
	// returns an error if the connection can not be established
	Connect(connString string) error

	// SSHConnect opens a ssh tunnel to the host and connects to a postgresql database there
	// taken from: https://github.com/jackc/pgx/issues/382
	SSHConnect(connString string, sshUser string, sshPassword string, sshHost string, sshPort int) error

	// Close stops the connection
	Close()

	// Ping executes a simple query against the database to check if the connection is healthy
	// returns the database version or an error
	Ping() (string, error)

	// query is the generic method to query the database
	// param receiver must be a pointer to a slice of struct that contains the expected columns as fields
	// param args can be a number of arguments to the query
	// returns any error occurring while executing the query
	//
	// Example:
	// sql := "SELECT phonenumber, name FROM phonebook WHERE name = '$1'" // use $i to fill the ith arg in the sql
	// receiver := []struct{ Phonenumber int, Name string }{} // be sure to use uppercase field names; make it a slice of your type because call to pgx.QueryRow will return a list of rows even if there is just one
	// err := query(sql, receiver, "Turing")
	Query(sql string, receiver interface{}, args ...interface{}) error
}

type postgresConnector struct {
	connection *pgxpool.Pool
}

// NewPostgresConnector returns a pointer to a new PostgresConnector instance
func NewPostgresConnector() PostgresConnector {
	return &postgresConnector{}
}

func (pC *postgresConnector) Connect(connString string) error {
	log.Info("Connecting to database...")
	connection, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return err
	}
	pC.connection = connection
	return nil
}

func (pC *postgresConnector) SSHConnect(dbHost string, dbPort int, dbUser string, dbPassword string, dbName string,
	sshUser string, sshPassword string, sshHost string, sshPort int) error {
	// The client configuration with configuration option to use the ssh-agent
	sshConfig := &ssh.ClientConfig{
		User:            sshUser,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // should be alright with the known ssh server
	}

	// When there's a non empty password add the password AuthMethod
	if sshPassword != "" {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PasswordCallback(func() (string, error) {
			return sshPassword, nil
		}))
	}

	sshcon, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshHost, sshPort), sshConfig)
	if err != nil {
		return fmt.Errorf("Can not connect to database via ssh: %v", err)
	}
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     dbHost,
			User:     dbUser,
			Password: dbPassword,
			Database: dbName,
			Dial: func(network, addr string) (net.Conn, error) {
				return sshcon.Dial(network, addr)
			},
		},
	}
	log.Info("Connecting to database via ssh...")
	connection, err := pgx.ConnectConfig(context.Background(), connPoolConfig)
	if err != nil {
		return fmt.Errorf("Can not create new connection pool: %v", err)
	}
	pC.connection = connection
	return nil
}

func (pC *postgresConnector) Close() {
	pC.connection.Close()
}

func (pC *postgresConnector) Ping() (string, error) {
	version := []struct{ Version string }{}
	err := pC.Query("SELECT version()", &version)
	if err != nil {
		return "", err
	}
	return version[0].Version, nil
}

func (pC *postgresConnector) Query(sql string, receiver interface{}, args ...interface{}) error {
	// from https://github.com/jackc/pgx/issues/878
	// Add PostgreSQL magic json functions
	// This gives us a single row back even if the query returns many rows
	// Now they get aggregated into a jsonb
	completeSql := fmt.Sprintf(
		`WITH orig_sql AS 
		(%s) 
		SELECT jsonb_agg(row_to_json(orig_sql.*)) 
		FROM orig_sql;`,
		sql)
	err := pC.connection.QueryRow(context.Background(), completeSql, args...).Scan(receiver)
	if err != nil {
		return fmt.Errorf("Can not query database: %w", err)
	}
	return nil
}
