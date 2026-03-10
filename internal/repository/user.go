package repository

import (
	"context"
	"time"

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
		return nil, ErrInvalidAPIKey
	}

	return &user, err
}

func (u *UserRepo) SaveRefreshToken(ctx context.Context, userID string, token string, expiresAt time.Time) error {

	_, err := u.pool.Exec(ctx, "INSERT INTO refresh_tokens(user_id, token, expires_at) VALUES ($1, $2, $3)", userID, token, expiresAt)

	if err != nil {
		return err
	}

	return nil

}

func (u *UserRepo) FindByRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken

	err := u.pool.QueryRow(ctx,
		"SELECT id, user_id, token, expires_at, is_revoked FROM refresh_tokens WHERE token=$1", token).
		Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.IsRevoked)

	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &rt, nil
}

func (u *UserRepo) RevokeRefreshToken(ctx context.Context, token string) error {
	_, err := u.pool.Exec(ctx, "UPDATE refresh_tokens SET is_revoked=true WHERE token=$1", token)
	return err
}
