package postgres

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	CantStartTransacMsg   = "can't start transaction"
	FailedToScanRowsMsg   = "failed to scan rows"
	FailedToScanRowMsg    = "failed to scan row"
	FailedToCommitTransac = "failed to commit transaction"
	FailedToStartQuery    = "failed to start query"
	FailedToUpdateValues  = "failed to update values"
)

const (
	TeamsTable     = "teams"
	UsersTable     = "users"
	PRTable        = "pull_requests"
	ReviewersTable = "pr_reviewers"
)

type Storage struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func New(ctx context.Context, storagePath string, log *slog.Logger) (*Storage, error) {

	pool, err := pgxpool.New(ctx, storagePath)
	if err != nil {
		log.Debug("couldn't create pgx pool", slog.Any("errror", err))
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Debug("ping is not successful", slog.Any("errror", err))
		return nil, err

	}

	return &Storage{
		pool: pool,
		log:  log,
	}, nil
}

func (s *Storage) Close() {
	s.pool.Close()
}
