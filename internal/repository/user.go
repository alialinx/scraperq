package repository

import (
	"context"
	"github.com/alialin/scraperq/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {

	return &UserRepo{pool: pool}

}

func (u *UserRepo) Create(ctx context.Context, user *models.User) error {

	var exists bool

	err := u.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 from users WHERE email = $1)", user.Email).Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		return ErrEmailExists
	}

	err = u.pool.QueryRow(ctx, "INSERT INTO users (email, password_hash, api_key, daily_limit, monthly_limit) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at",
		user.Email, user.PasswordHash, user.APIKey, user.DailyLimit, user.MonthlyLimit).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {

	var user models.User

	err := u.pool.QueryRow(ctx, "SELECT id, email, password_hash, api_key, is_active, daily_limit, monthly_limit, created_at FROM users WHERE email=$1", email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.APIKey, &user.IsActive, &user.DailyLimit, &user.MonthlyLimit, &user.CreatedAt)

	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	}

	return &user, err
}

func (u *UserRepo) FindByAPIKey(ctx context.Context, apiKey string) (*models.User, error) {

	var user models.User

	err := u.pool.QueryRow(ctx, "SELECT id, email, password_hash, api_key, is_active, daily_limit, monthly_limit, created_at FROM users WHERE api_key=$1", apiKey).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.APIKey, &user.IsActive, &user.DailyLimit, &user.MonthlyLimit, &user.CreatedAt)

	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	}

	return &user, err
}
