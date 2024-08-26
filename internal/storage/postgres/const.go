package postgres

import "time"

const (
	defaultMaxOpenConnections    = 20
	defaultMaxIdleConnections    = 1
	defaultMaxConnectionLifetime = 300 * time.Second
	defaultMaxConnectionIdleTime = 60 * time.Second
)
