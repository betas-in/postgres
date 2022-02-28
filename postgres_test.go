package postgres

import (
	"database/sql"
	"testing"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
)

func TestPostgres(t *testing.T) {
	conf := Config{}
	conf.Host = "localhost"
	conf.Port = 7005
	conf.Username = "postgres"
	conf.Password = "596a96cc7bf9108cd896f33c44aedc8a"
	conf.DatabaseName = "testpg"
	log := logger.NewLogger(3, true)

	d, err := NewDatabase(&conf, log)
	utils.Test().Nil(t, err)

	query := "SELECT * FROM pg_catalog.pg_tables LIMIT 1"
	tables := []struct {
		Schemaname  string
		Tablename   string
		Tableowner  string
		Tablespace  sql.NullString
		Hasindexes  bool
		Hasrules    bool
		Hastriggers bool
		Rowsecurity bool
	}{}
	err = d.Select(&tables, query)
	utils.Test().Nil(t, err)

	err = d.Close()
	utils.Test().Nil(t, err)
}
