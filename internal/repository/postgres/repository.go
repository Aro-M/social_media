package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"social_media/internal/config"
	"social_media/internal/entity"
	"social_media/internal/repository"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository() (repository.Repository, error) {
	db, err := sql.Open("postgres", config.DSN())
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout())
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &Repository{db: db}, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) CreateUser(ctx context.Context, u *entity.User) (*entity.User, error) {
	query := `INSERT INTO users (id, first_name, last_name, email, age, password_hash, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, u.ID, u.FirstName, u.LastName, u.Email, u.Age, u.PasswordHash, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("exec context: %w", err)
	}
	return u, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var u entity.User
	query := `SELECT id, first_name, last_name, email, age, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Age, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("query row context: %w", err)
	}
	return &u, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var u entity.User
	query := `SELECT id, first_name, last_name, email, age, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Age, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("query row context: %w", err)
	}
	return &u, nil
}

func (r *Repository) SearchUsers(ctx context.Context, firstName, lastName string, age int) ([]*entity.User, error) {
	// Simple mock-like implementation for now as it's a refactoring task
	return nil, fmt.Errorf("not implemented")
}

func (r *Repository) CreateFriendRequest(ctx context.Context, fr *entity.Friend) (*entity.Friend, error) {
	query := `INSERT INTO friends (id, sender_id, receiver_id, status, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, fr.ID, fr.SenderID, fr.ReceiverID, fr.Status, fr.CreatedAt, fr.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("exec context: %w", err)
	}
	return fr, nil
}

func (r *Repository) GetFriendRequest(ctx context.Context, id uuid.UUID) (*entity.Friend, error) {
	var f entity.Friend
	query := `SELECT id, sender_id, receiver_id, status, created_at, updated_at FROM friends WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&f.ID, &f.SenderID, &f.ReceiverID, &f.Status, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("query row context: %w", err)
	}
	return &f, nil
}

func (r *Repository) UpdateFriendRequest(ctx context.Context, id uuid.UUID, status entity.FriendStatus) error {
	query := `UPDATE friends SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("exec context: %w", err)
	}
	return nil
}

func (r *Repository) ListFriendRequests(ctx context.Context, userID uuid.UUID) ([]*entity.Friend, error) {
	return nil, fmt.Errorf("not implemented")
}
