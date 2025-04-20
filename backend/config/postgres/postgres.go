package postgres

import (
	"pvz/internal/utils"
)

const (
	defaultDataBaseURL string = "postgresql://quickflow_admin:SuperSecurePassword1@localhost:5432/quickflow_db"
)

type PostgresConfig struct {
	dataBaseURL string
}

func NewPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		dataBaseURL: utils.GetEnv("DATABASE_URL", defaultDataBaseURL),
	}
}

func (p *PostgresConfig) GetURL() string {
	return p.dataBaseURL
}
