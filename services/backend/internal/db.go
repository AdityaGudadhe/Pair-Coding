package internal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultMaxConns          = 25
	defaultMinConns          = 1
	defaultMaxConnLifetime   = time.Hour
	defaultMaxConnIdleTime   = 30 * time.Minute
	defaultHealthCheckPeriod = time.Minute
	defaultConnectTimeout    = 5 * time.Second
)

// InitDatabase builds a pgx connection pool from DATABASE_URL, applies pool
// tuning (overridable via DB_* env vars), and verifies connectivity with a ping.
// The caller owns the returned pool and must Close it on shutdown.
func InitDatabase(ctx context.Context) (*pgxpool.Pool, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
	}

	maxConns, err := envInt("DB_MAX_CONNS", defaultMaxConns)
	if err != nil {
		return nil, err
	}
	minConns, err := envInt("DB_MIN_CONNS", defaultMinConns)
	if err != nil {
		return nil, err
	}
	if minConns > maxConns {
		return nil, fmt.Errorf("DB_MIN_CONNS (%d) cannot exceed DB_MAX_CONNS (%d)", minConns, maxConns)
	}
	maxConnLifetime, err := envDuration("DB_MAX_CONN_LIFETIME", defaultMaxConnLifetime)
	if err != nil {
		return nil, err
	}
	maxConnIdleTime, err := envDuration("DB_MAX_CONN_IDLE_TIME", defaultMaxConnIdleTime)
	if err != nil {
		return nil, err
	}
	healthCheckPeriod, err := envDuration("DB_HEALTH_CHECK_PERIOD", defaultHealthCheckPeriod)
	if err != nil {
		return nil, err
	}
	connectTimeout, err := envDuration("DB_CONNECT_TIMEOUT", defaultConnectTimeout)
	if err != nil {
		return nil, err
	}

	config.MaxConns = int32(maxConns)
	config.MinConns = int32(minConns)
	config.MaxConnLifetime = maxConnLifetime
	config.MaxConnLifetimeJitter = maxConnLifetime / 10
	config.MaxConnIdleTime = maxConnIdleTime
	config.HealthCheckPeriod = healthCheckPeriod
	config.ConnConfig.ConnectTimeout = connectTimeout
	if _, ok := config.ConnConfig.RuntimeParams["application_name"]; !ok {
		config.ConnConfig.RuntimeParams["application_name"] = "pair-coding-backend"
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

func envInt(key string, fallback int) (int, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback, nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v < 0 {
		return 0, fmt.Errorf("%s must be a non-negative integer, got %q", key, raw)
	}
	return v, nil
}

func envDuration(key string, fallback time.Duration) (time.Duration, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback, nil
	}
	v, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be a duration like 30s or 5m, got %q", key, raw)
	}
	return v, nil
}
