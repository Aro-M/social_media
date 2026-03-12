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

// NewRepository connects to the Postgres database.
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

// Close shuts down the database connection.
func (r *Repository) Close() error {
	return r.db.Close()
}

// CreateUser inserts a new user record.
func (r *Repository) CreateUser(ctx context.Context, u *entity.User) (*entity.User, error) {
	query := `INSERT INTO users (id, first_name, last_name, email, age, password_hash, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, u.ID, u.FirstName, u.LastName, u.Email, u.Age, u.PasswordHash, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("exec context: %w", err)
	}
	return u, nil
}

// GetUserByID finds a user by their ID.
func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var u entity.User
	query := `SELECT id, first_name, last_name, email, age, password_hash, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Age, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("query row context: %w", err)
	}
	return &u, nil
}

// GetUserByEmail finds a user by their email.
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var u entity.User
	query := `SELECT id, first_name, last_name, email, age, password_hash, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Age, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("query row context: %w", err)
	}
	return &u, nil
}

// UpdateUser saves changed user information.
func (r *Repository) UpdateUser(ctx context.Context, u *entity.User) error {
	query := `UPDATE users SET first_name = $1, last_name = $2, email = $3, age = $4, password_hash = $5, updated_at = $6 WHERE id = $7`
	res, err := r.db.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.Age, u.PasswordHash, u.UpdatedAt, u.ID)
	if err != nil {
		return fmt.Errorf("exec context: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// SearchUsers finds users by name or age, excluding the current user.
func (r *Repository) SearchUsers(ctx context.Context, firstName, lastName string, age int, currentUserID uuid.UUID) ([]*entity.User, error) {
	query := `SELECT id, first_name, last_name, email, age, created_at, updated_at FROM users WHERE id != $1`
	args := []interface{}{currentUserID}
	argID := 2

	if firstName != "" {
		query += fmt.Sprintf(" AND first_name ILIKE $%d", argID)
		args = append(args, firstName+"%")
		argID++
	}

	if lastName != "" {
		query += fmt.Sprintf(" AND last_name ILIKE $%d", argID)
		args = append(args, lastName+"%")
		argID++
	}

	if age > 0 {
		query += fmt.Sprintf(" AND age = $%d", argID)
		args = append(args, age)
		argID++
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Age, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return users, nil
}

// CreateFriendRequest inserts a new friend request.
func (r *Repository) CreateFriendRequest(ctx context.Context, fr *entity.Friend) (*entity.Friend, error) {
	query := `INSERT INTO friends (id, sender_id, receiver_id, status, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, fr.ID, fr.SenderID, fr.ReceiverID, fr.Status, fr.CreatedAt, fr.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("exec context: %w", err)
	}
	return fr, nil
}

// GetFriendRequest finds a request by its ID.
func (r *Repository) GetFriendRequest(ctx context.Context, id uuid.UUID) (*entity.Friend, error) {
	var f entity.Friend
	query := `SELECT id, sender_id, receiver_id, status, created_at, updated_at FROM friends WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&f.ID, &f.SenderID, &f.ReceiverID, &f.Status, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("query row context: %w", err)
	}
	return &f, nil
}

// UpdateFriendRequest updates the request status.
func (r *Repository) UpdateFriendRequest(ctx context.Context, id uuid.UUID, status entity.FriendStatus) error {
	query := `UPDATE friends SET status = $1, updated_at = NOW() WHERE id = $2`
	res, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("exec context: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("friend request not found or not updated")
	}

	return nil
}

// ListFriendRequests gets all pending user requests.
func (r *Repository) ListFriendRequests(ctx context.Context, userID uuid.UUID) ([]*entity.Friend, error) {
	query := `SELECT id, sender_id, receiver_id, status, created_at, updated_at 
	          FROM friends 
			  WHERE receiver_id = $1 AND status = 'pending' 
			  ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}
	defer rows.Close()

	var requests []*entity.Friend
	for rows.Next() {
		var fr entity.Friend
		if err := rows.Scan(&fr.ID, &fr.SenderID, &fr.ReceiverID, &fr.Status, &fr.CreatedAt, &fr.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		requests = append(requests, &fr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	if requests == nil {
		requests = make([]*entity.Friend, 0)
	}

	return requests, nil
}

// GetFriendshipByParticipants finds any existing friend record between two users.
func (r *Repository) GetFriendshipByParticipants(ctx context.Context, userID1, userID2 uuid.UUID) (*entity.Friend, error) {
	var f entity.Friend
	query := `SELECT id, sender_id, receiver_id, status, created_at, updated_at 
	          FROM friends 
	          WHERE (sender_id = $1 AND receiver_id = $2) 
	             OR (sender_id = $2 AND receiver_id = $1)`
	err := r.db.QueryRowContext(ctx, query, userID1, userID2).Scan(&f.ID, &f.SenderID, &f.ReceiverID, &f.Status, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found is not an error here
		}
		return nil, fmt.Errorf("query row context: %w", err)
	}
	return &f, nil
}

// ListFriends finds all accepted friends.
func (r *Repository) ListFriends(ctx context.Context, userID uuid.UUID) ([]*entity.User, error) {
	query := `
		SELECT u.id, u.first_name, u.last_name, u.email, u.age, u.created_at, u.updated_at
		FROM users u
		JOIN friends f ON (
			(f.sender_id = $1 AND f.receiver_id = u.id)
			OR (f.receiver_id = $1 AND f.sender_id = u.id)
		)
		WHERE f.status = 'accepted'
		ORDER BY u.first_name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}
	defer rows.Close()

	var friends []*entity.User
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Age, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		friends = append(friends, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	if friends == nil {
		friends = make([]*entity.User, 0)
	}

	return friends, nil
}

// RemoveFriend deletes a friend record.
func (r *Repository) RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	query := `DELETE FROM friends 
	          WHERE ((sender_id = $1 AND receiver_id = $2) 
			     OR (sender_id = $2 AND receiver_id = $1))`

	res, err := r.db.ExecContext(ctx, query, userID, friendID)
	if err != nil {
		return fmt.Errorf("exec context: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("friendship not found")
	}

	return nil
}
