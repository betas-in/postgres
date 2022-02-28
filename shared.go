package postgres

import (
	"fmt"
)

// ConnectionString - Create a connection string
func ConnectionString(conf *Config) string {
	return fmt.Sprintf("dbname=%s user=%s password='%s' host=%s port=%d sslmode=%s", conf.DatabaseName, conf.Username, conf.Password, conf.Host, conf.Port, "disable")
}

// ConnectionURL - Create a connection URL
func ConnectionURL(conf *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&timezone=UTC", conf.Username, conf.Password, conf.Host, conf.Port, conf.DatabaseName)
}
