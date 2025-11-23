package postgres

import (
	"avito_test/internal/domain/models"
	"avito_test/internal/storage"
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) CreateTeam(ctx context.Context, team *models.Team) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.log.Debug(CantStartTransacMsg, slog.Any("error", err))
		return err
	}

	defer tx.Rollback(ctx)
	query := fmt.Sprintf("INSERT INTO %s (team_name) VALUES ($1)", TeamsTable)
	_, err = tx.Exec(ctx, query, team.TeamName)
	if err != nil {
		s.log.Debug("can't insert teamname", slog.Any("error", err), slog.String("team_name", team.TeamName))
		return storage.ErrTeamAlreadyExists
	}
	batch := &pgx.Batch{}
	query = fmt.Sprintf("INSERT INTO %s (user_id, username, team_name, is_active) VALUES ($1, $2, $3, $4) ON CONFLICT (user_id) DO UPDATE SET team_name = $3, is_active = $4", UsersTable)
	for _, m := range team.Members {
		batch.Queue(query, m.UserID, m.Username, team.TeamName, m.IsActive)
	}
	br := tx.SendBatch(ctx, batch)
	for range team.Members {
		if _, err := br.Exec(); err != nil {
			s.log.Debug("failed to read batch results", slog.Any("error", err))
			return err
		}
	}
	br.Close()
	err = tx.Commit(ctx)
	if err != nil {
		s.log.Warn(FailedToCommitTransac, slog.Any("error", err))
		return err
	}
	return nil
}

func (s *Storage) GetTeamWithMembers(ctx context.Context, teamName string) (*models.Team, error) {

	team := models.NewTeam(teamName)
	query := fmt.Sprintf("SELECT user_id, username, is_active FROM %s WHERE team_name = $1", UsersTable)
	rows, err := s.pool.Query(ctx, query, teamName)
	if err != nil {
		s.log.Debug(FailedToStartQuery, slog.Any("error", err), slog.String("team_name", teamName))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var member models.TeamMember
		if err := rows.Scan(&member.UserID, &member.Username, &member.IsActive); err != nil {
			s.log.Debug(FailedToScanRowsMsg, slog.Any("error", err))
			return nil, err
		}
		team.Members = append(team.Members, &member)
	}
	if err := rows.Err(); err != nil {
		s.log.Debug("error after reading", slog.Any("error", err))
		return nil, err
	}
	if len(team.Members) == 0 {
		return nil, storage.ErrTeamNotFound
	}

	return team, nil
}
