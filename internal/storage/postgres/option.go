package postgres

import "time"

type Option func(*Postgres)

func MaxOpenConnections(size int) Option {
	return func(c *Postgres) {
		c.maxOpenConnections = size
	}
}

func MaxIdleConnections(c int) Option {
	return func(p *Postgres) {
		p.maxIdleConnections = c
	}
}

func MaxConnectionLifeTime(d time.Duration) Option {
	return func(p *Postgres) {
		p.maxConnectionLifeTime = d
	}
}

func MaxConnectionIdleTime(d time.Duration) Option {
	return func(p *Postgres) {
		p.maxConnectionIdleTime = d
	}
}
