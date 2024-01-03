package repository

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
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

	Connection() *pgxpool.Pool

	// Ping executes a simple Ping against the database to check if the connection is healthy
	// returns an error if the Ping failed
	Ping() error

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
	connection, err := pgxpool.New(context.Background(), connString)
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
		host, portString, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		port, err := strconv.Atoi(portString)
		if err != nil {
			return nil, err
		}
		// with sshcon.Dial() the remoteAddr is empty (0.0.0.0) so no CancelRequest could be sent.
		return sshcon.DialTCP(network, nil, &net.TCPAddr{
			IP:   net.ParseIP(host),
			Port: port,
		})
	}
	log.Info("Connecting to database via ssh...")
	connection, err := pgxpool.NewWithConfig(context.Background(), connPoolConfig)
	if err != nil {
		return fmt.Errorf("Can not create new connection pool: %v", err)
	}
	pC.connection = connection
	return nil
}

func (pC *postgresConnector) Close() {
	pC.connection.Close()
}

func (pC *postgresConnector) Connection() *pgxpool.Pool {
	return pC.connection
}

func (pC *postgresConnector) Ping() error {
	c, err := pC.connection.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer c.Release()
	return c.Ping(context.Background())
}

func (pC *postgresConnector) Query(ctx context.Context, sql string, receiver interface{}, args ...interface{}) error {
	// from https://github.com/jackc/pgx/issues/878
	// Add PostgreSQL magic json functions
	// This gives us a single row back even if the query returns many rows
	// Now they get aggregated into a jsonb
	// Note: jsonb has a maximum size. This is exceeded with the additional fields (rocktype & rockclass) in the samples query so we need another solution
	completeSql := fmt.Sprintf(
		`WITH orig_sql AS 
		(%s) 
		SELECT jsonb_agg(row_to_json(orig_sql.*)) 
		FROM orig_sql;`,
		sql)

	// manually acquire and release connection to be able to send CancelRequest() on context canceled by client
	c, err := pC.connection.Acquire(ctx)
	if err != nil {
		return err
	}
	defer c.Release()
	stopChan := make(chan bool)
	go cancelQueryOnContextCanceled(ctx, c, stopChan)
	row := c.QueryRow(ctx, completeSql, args...)
	err = row.Scan(receiver)
	stopChan <- true
	if err != nil {
		return fmt.Errorf("Can not query database: %w", err)
	}
	return nil
}

func Query[T any](ctx context.Context, pC PostgresConnector, sql string, args ...interface{}) ([]T, error) {
	// manually acquire and release connection to be able to send CancelRequest() on context canceled by client
	c, err := pC.Connection().Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Release()
	stopChan := make(chan bool)
	go cancelQueryOnContextCanceled(ctx, c, stopChan)
	rows, err := c.Query(ctx, sql, args...)
	stopChan <- true
	if err != nil {
		return nil, fmt.Errorf("Can not query database: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[T])
}

// cancelQueryOnContextCanceled is an async context watcher to send a CancelRequest() if the context is canceled by client
func cancelQueryOnContextCanceled(ctx context.Context, c *pgxpool.Conn, stopChan chan bool) {
	// block until context is done or query returned
	select {
	case <-stopChan:
		// stop goroutine to avoid call to c.Conn() after c.Release()
		return
	case <-ctx.Done():
		err := c.Conn().PgConn().CancelRequest(context.Background())
		if err != nil {
			// query cancellation failed
			log.Warn(fmt.Sprintf("CancelRequest failed: %s", err.Error()))
		}
	}
}
