package postgres

import (
	"barbot/internal/config"
	"database/sql"
	"fmt"
)

func New(cfg *config.DatabaseConfig) (*sql.DB, error) {
	var connectString = fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.UserName,
		cfg.Password,
		cfg.DbName,
	)
	db, err := sql.Open("postgres", connectString)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}
