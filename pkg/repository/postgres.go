package repository

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v4/pgxpool"
)

// PostgresConnector interface exposes methods to connect to and interact with a postgreSQL instance
type PostgresConnector interface {
	// Connect tries to connect to a postgreSQL database with the credentials provided in the connString
	// returns an error if the connection can not be established
	Connect(connString string) error

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
