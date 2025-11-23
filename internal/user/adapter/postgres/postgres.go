package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/silentmol/avito-backend-trainee/internal/apperr"
	"github.com/silentmol/avito-backend-trainee/internal/user/domain"
)

type UserRepository struct {
	conn *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) *UserRepository {
	return &UserRepository{conn: conn}
}

func (u *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
	INSERT INTO users (id, name, team_name, is_active)
	VALUES ($1, $2, $3, $4)
	RETURNING id, name, team_name, is_active
	`

	var createdUser domain.User
	err := u.conn.QueryRow(ctx, query,
		user.ID,
		user.Name,
		user.TeamName,
		user.IsActive,
	).Scan(
		&createdUser.ID,
		&createdUser.Name,
		&createdUser.TeamName,
		&createdUser.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("db: failed to create user: %w", err)
	}

	return &createdUser, nil
}

func (u *UserRepository) GetUser(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User

	query := `
		SELECT id, name, team_name, is_active 
		FROM users 
		WHERE id=$1
	`
	err := u.conn.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.TeamName,
		&user.IsActive,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, fmt.Errorf("db: failed to get user: %w", err)
	}
	return &user, nil
}

func (u *UserRepository) SetIsActive(ctx context.Context, id string, isActive bool) (*domain.User, error) {
	var user domain.User

	query := `
		UPDATE users
		SET is_active = $1
		WHERE id = $2
		RETURNING id, name, team_name, is_active
	`

	err := u.conn.QueryRow(ctx, query, isActive, id).Scan(
		&user.ID,
		&user.Name,
		&user.TeamName,
		&user.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, fmt.Errorf("db: failed to update is_active: %w", err)
	}

	return &user, nil
}
