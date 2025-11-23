package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/silentmol/avito-backend-trainee/internal/apperr"
	"github.com/silentmol/avito-backend-trainee/internal/team/domain"
)

type TeamRepository struct {
	conn *pgxpool.Pool
}

func NewTeamRepository(conn *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{conn: conn}
}

func (t *TeamRepository) CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	insertTeamQuery := `
		INSERT INTO teams (name)
		VALUES ($1)
		RETURNING name
	`

	var createdTeam domain.Team
	if err := t.conn.QueryRow(ctx, insertTeamQuery, team.Name).Scan(&createdTeam.Name); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, apperr.ErrTeamExists
		}
		return nil, fmt.Errorf("db: failed to create team: %w", err)
	}

	upsertUserQuery := `
		INSERT INTO users (id, name, team_name, is_active)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name,
		    team_name = EXCLUDED.team_name,
		    is_active = EXCLUDED.is_active
	`

	for _, member := range team.Members {
		if _, err := t.conn.Exec(
			ctx,
			upsertUserQuery,
			member.ID,
			member.Name,
			team.Name,
			member.IsActive,
		); err != nil {
			return nil, fmt.Errorf("db: failed to upsert user %s: %w", member.ID, err)
		}
	}

	createdTeam.Members = team.Members

	return &createdTeam, nil
}

func (t *TeamRepository) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	getTeamQuery := `
		SELECT name
		FROM teams
		WHERE name = $1
	`

	var name string
	if err := t.conn.QueryRow(ctx, getTeamQuery, teamName).Scan(&name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, fmt.Errorf("db: failed to get team: %w", err)
	}

	getMembersQuery := `
		SELECT id, name, is_active
		FROM users
		WHERE team_name = $1
	`

	rows, err := t.conn.Query(ctx, getMembersQuery, teamName)
	if err != nil {
		return nil, fmt.Errorf("db: failed to get team members: %w", err)
	}
	defer rows.Close()

	team := &domain.Team{
		Name:    name,
		Members: make([]domain.TeamMember, 0),
	}

	for rows.Next() {
		var member domain.TeamMember
		if err := rows.Scan(&member.ID, &member.Name, &member.IsActive); err != nil {
			return nil, fmt.Errorf("db: failed to scan team member: %w", err)
		}
		team.Members = append(team.Members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: rows error: %w", err)
	}

	return team, nil
}
