package main

import (
	"context"
	"database/sql"
	"time"
)

func openDB(c config) (*sql.DB, error) {
	db, err := sql.Open("postgres", c.db.dsn)
	if err != nil {
		return nil, err
	}

	// Set maximum open connections
	db.SetMaxOpenConns(c.db.maxOpenConns)
	// Set maximum idle connections
	db.SetMaxIdleConns(c.db.maxIdleConns)

	// Convert maxIdleTime of type string to time.Duration
	duration, err := time.ParseDuration(c.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	// Set the maximum idle timeout
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
