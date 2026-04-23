package repository

// ConnectionParams holds the parameters for a database connection
type ConnectionParams struct {
	Host        string
	Port        int
	Username    string
	Password    string
	DBName      string
	SSHHost     string
	SSHPort     int
	SSHUser     string
	SSHPassword string
}
