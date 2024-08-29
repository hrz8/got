package httpserver

import "time"

const (
	// app
	defaultPort            = 5101
	defaultShutdownTimeout = 15 * time.Second
	// standard
	defaultReadHeaderTimeout = 5 * time.Second
	defaultReadTimeout       = 10 * time.Second
	defaultWriteTimeout      = 15 * time.Second
	defaultIdleTimeout       = 180 * time.Second
)
