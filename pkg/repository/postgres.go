package repository

import (
	"context"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"

	"github.com/jackc/pgx/v4/pgxpool"
)

// PostgresConnector interface exposes methods to connect to and interact with a postgreSQL instance
type PostgresConnector interface {
	// Connect tries to connect to a postgreSQL database with the credentials provided in the connString
	// returns an error if the connection can not be established
	Connect(connString string) error

	// SSHConnect opens a ssh tunnel to the host and connects to a postgresql database there
	// taken from: https://github.com/jackc/pgx/issues/382
	SSHConnect(connString string, params *ConnectionParams) error

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
	Query(ctx context.Context, sql string, receiver interface{}, args ...interface{}) error
}

type postgresConnector struct {
	connection *pgxpool.Pool
}

// ConnectionParams holds the parameters for a postgresql database connection
type ConnectionParams struct {
	DBHost      string
	DBPort      int
	DBUser      string
	DBPassword  string
	DBName      string
	SSHHost     string
	SSHPort     int
	SSHUser     string
	SSHPassword string
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

func (pC *postgresConnector) SSHConnect(connString string, params *ConnectionParams) error {
	// The client configuration with configuration option to use the ssh-agent
	sshConfig := &ssh.ClientConfig{
		User:            params.SSHUser,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // should be alright with the known ssh server
	}

	// When there's a non empty password add the password AuthMethod
	if params.SSHPassword != "" {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PasswordCallback(func() (string, error) {
			return params.SSHPassword, nil
		}))
	}

	sshcon, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", params.SSHHost, params.SSHPort), sshConfig)
	if err != nil {
		return fmt.Errorf("Can not connect to database via ssh: %v", err)
	}

	connPoolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return fmt.Errorf("Can not parse config: %v", err)
	}
	connPoolConfig.ConnConfig.DialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return sshcon.Dial(network, addr)
	}
	log.Info("Connecting to database via ssh...")
	connection, err := pgxpool.ConnectConfig(context.Background(), connPoolConfig)
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
	err := pC.Query(context.Background(), "SELECT version()", &version)
	if err != nil {
		return "", err
	}
	return version[0].Version, nil
}

func (pC *postgresConnector) Query(ctx context.Context, sql string, receiver interface{}, args ...interface{}) error {
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
	err := pC.connection.QueryRow(ctx, completeSql, args...).Scan(receiver)
	if err != nil {
		return fmt.Errorf("Can not query database: %w", err)
	}
	return nil
}
