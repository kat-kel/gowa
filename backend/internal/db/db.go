package db

import (
	sql "database/sql"
)

type Store struct {
	DB *sql.DB
}

func Open(dsn string) (*Store, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Store{DB: db}, nil
}

func (s *Store) Migrate() error {
	_, err := s.DB.Exec(`
	CREATE TABLE IF NOT EXISTS actors (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		voice TEXT NOT NULL
	)`)
	return err
}
