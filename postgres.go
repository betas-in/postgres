package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/betas-in/logger"
	"github.com/jmoiron/sqlx"
)

var (
	// ErrDBNotDefined ...
	ErrDBNotDefined = errors.New("Database connection is not defined yet")
)

// Database is the External Database definition
type Database interface {
	GetDB() *sqlx.DB
	Select(dest interface{}, query string, args ...interface{}) error
	// Get(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRow(query string, args ...interface{}) *sqlx.Row
	// QueryBind(argCount, rowCount int) string
	// Rebind(query string) string
	Close() error
}

// NewDatabase - Create a new Database
func NewDatabase(conf *Config, logger *logger.Logger) (Database, error) {
	pg := postgres{conf: conf, log: logger}
	pg.defaults()
	err := pg.Connect()
	if err != nil {
		logger.Error("postgres.new").Msgf("Error in connecting to database: %+v", err)
		return &pg, err
	}
	return &pg, nil
}

// Internal Postgres implementation
type postgres struct {
	db   *sqlx.DB
	conf *Config
	log  *logger.Logger
}

type Config struct {
	Host               string
	Port               int
	DatabaseName       string
	Username           string
	Password           string
	MaxIdleConnections int
	MaxOpenConnections int
	MaxLifetime        int
	MigrationPath      string
}

func (p *postgres) defaults() {
	if p.conf.Host == "" {
		p.conf.Host = "localhost"
	}
	if p.conf.Port == 0 {
		p.conf.Port = 5432
	}
	if p.conf.Username == "" {
		p.conf.Username = "postgres"
	}
	if p.conf.MaxIdleConnections == 0 {
		p.conf.MaxIdleConnections = 20
	}
	if p.conf.MaxOpenConnections == 0 {
		p.conf.MaxOpenConnections = 20
	}
	if p.conf.MaxLifetime == 0 {
		p.conf.MaxLifetime = 5
	}
}

// Load the database
func (p *postgres) Connect() error {
	var err error

	p.db, err = sqlx.Open("postgres", ConnectionString(p.conf))
	if err != nil {
		p.log.Error("postgres.connect").Msgf("Could not open db %+v", ConnectionString(p.conf))
		return err
	}

	err = p.db.Ping()
	if err != nil {
		p.log.Error("postgres.connect").Msgf("Could not ping db %+v", err)
		return err
	}

	p.db.SetMaxIdleConns(p.conf.MaxIdleConnections)
	p.db.SetMaxOpenConns(p.conf.MaxOpenConnections)
	p.db.SetConnMaxLifetime(time.Duration(p.conf.MaxLifetime) * time.Second)

	p.log.Info("postgres.connect").Msgf("connected to %s @ %s:%d", p.conf.DatabaseName, p.conf.Host, p.conf.Port)
	return nil
}

// GetDB the database object
func (p *postgres) GetDB() *sqlx.DB {
	return p.db
}

// Select data from DB
func (p *postgres) Select(dest interface{}, query string, args ...interface{}) error {
	if p.db == nil {
		return ErrDBNotDefined
	}

	err := p.db.Select(dest, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// Get data from DB
func (p *postgres) Get(dest interface{}, query string, args ...interface{}) error {
	err := p.db.Get(dest, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// Exec executes a query to DB
func (p *postgres) Exec(query string, args ...interface{}) (sql.Result, error) {
	if p.db == nil {
		return nil, ErrDBNotDefined
	}

	result, err := p.db.Exec(query, args...)
	if err != nil {
		return result, err
	}
	return result, err
}

// Query executes and returns rows
func (p *postgres) Query(query string, args ...interface{}) (*sqlx.Rows, error) {
	if p.db == nil {
		return nil, ErrDBNotDefined
	}

	rows, err := p.db.Queryx(query, args...)
	if err != nil {
		return rows, err
	}
	return rows, nil
}

func (p *postgres) QueryRow(query string, args ...interface{}) *sqlx.Row {
	if p.db == nil {
		return nil
	}
	return p.db.QueryRowx(query, args...)
}

func (p *postgres) Close() error {
	err := p.db.Close()
	if err != nil {
		p.log.Error("postgres.close").Msgf("Failed to close database: %+v", err)
	}
	p.db = nil
	return nil
}

// // Rebind function
// func (p *postgres) Rebind(query string) string {
// 	return p.db.Rebind(query)
// }
