package main

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	errUserNotFound    = errors.New("user not found")
	errProblemNotFound = errors.New("problem not found")
)

type appDatabase struct {
	conn *sql.DB
}

var database *appDatabase

func initDatabase() error {
	if database != nil {
		return errors.New("database is already initialized")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return err
	}

	database = &appDatabase{conn: db}
	return nil
}

func (db *appDatabase) Close() error {
	return db.conn.Close()
}

func (db *appDatabase) AdminUserIDByCredentials(ctx context.Context, user User) (int32, error) {
	var (
		userID int32
		err    error
	)

	err = db.conn.QueryRowContext(
		ctx,
		"SELECT user_id FROM users WHERE user_name = $1 AND password = $2",
		user.userName,
		user.password,
	).Scan(&userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errUserNotFound
		}

		return 0, err
	}
	return userID, nil
}

func (db *appDatabase) createNewProblem(ctx context.Context, problem Problem) (int32, error) {
	var problemID int32
	err := db.conn.QueryRowContext(
		ctx,
		`INSERT INTO problems (
			author_id,
			title,
			body,
			constraints,
			time_limit,
			memory_limit
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING problem_id`,
		problem.AuthorID,
		problem.Name,
		problem.Body,
		problem.Constraints,
		problem.TimeLimit,
		problem.MemoryLimit,
	).Scan(&problemID)
	if err != nil {
		return 0, err
	}

	return problemID, nil
}

func (db *appDatabase) getProblem(ctx context.Context, problemID int32) (Problem, error) {
	var problem Problem
	err := db.conn.QueryRowContext(
		ctx,
		`SELECT
			problem_id,
			author_id,
			title,
			body,
			constraints,
			time_limit,
			memory_limit
		FROM problems
		WHERE problem_id = $1`,
		problemID,
	).Scan(
		&problem.ProblemID,
		&problem.AuthorID,
		&problem.Name,
		&problem.Body,
		&problem.Constraints,
		&problem.TimeLimit,
		&problem.MemoryLimit,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Problem{}, errProblemNotFound
		}

		return Problem{}, err
	}

	return problem, nil
}

func (db *appDatabase) uploadInput(ctx context.Context, userID int32, problemID int32, inputID int32) error {
	_, err := db.conn.ExecContext(
		ctx,
		`INSERT INTO inputs (
			user_id,
			problem_id,
			input_storage_id
		) VALUES ($1, $2, $3)`,
		userID,
		problemID,
		inputID,
	)
	return err
}

func (db *appDatabase) uploadOutput(ctx context.Context, userID int32, problemID int32, outputID int32) error {
	_, err := db.conn.ExecContext(
		ctx,
		`INSERT INTO outputs (
			user_id,
			problem_id,
			output_storage_id
		) VALUES ($1, $2, $3)`,
		userID,
		problemID,
		outputID,
	)
	return err
}

func (db *appDatabase) uploadChecker(ctx context.Context, userID int32, problemID int32, checkerID int32) error {
	_, err := db.conn.ExecContext(
		ctx,
		`INSERT INTO checkers (
			user_id,
			problem_id,
			checker_storage_id
		) VALUES ($1, $2, $3)`,
		userID,
		problemID,
		checkerID,
	)
	return err
}
