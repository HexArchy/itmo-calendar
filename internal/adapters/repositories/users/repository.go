package users

import (
	"context"
	"fmt"
	"strings"

	"github.com/hexarchy/itmo-calendar/internal/entities"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, isu int64) (*entities.User, error) {
	const query = `
INSERT INTO users (isu, created_at, updated_at)
VALUES ($1, NOW(), NOW())
ON CONFLICT (isu) DO UPDATE SET updated_at = EXCLUDED.updated_at
RETURNING isu, created_at, updated_at
	`
	var user entities.User
	err := r.db.QueryRow(ctx, query, isu).Scan(&user.ISU, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "insert user")
	}

	return &user, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]entities.User, error) {
	const query = `
SELECT isu, created_at, updated_at
FROM users
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "select users")
	}
	defer rows.Close()

	var users []entities.User

	for rows.Next() {
		var u entities.User
		err = rows.Scan(&u.ISU, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "scan user")
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows error")
	}

	return users, nil
}

// GetBatch retrieves a batch of users with pagination.
func (r *Repository) GetBatch(ctx context.Context, limit, offset int) ([]entities.User, error) {
	const query = `
SELECT isu, created_at, updated_at
FROM users
ORDER BY isu
LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "select users batch")
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var u entities.User
		err = rows.Scan(&u.ISU, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "scan user")
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows error")
	}

	return users, nil
}

// FindByIDs retrieves users by their IDs.
func (r *Repository) FindByIDs(ctx context.Context, isus []int64) ([]entities.User, error) {
	if len(isus) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(isus))
	args := make([]any, len(isus))
	for i, u := range isus {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = u
	}

	query := `
SELECT isu, created_at, updated_at
FROM users
WHERE isu IN (` + strings.Join(placeholders, ",") + `)`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "find users by ids")
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var u entities.User
		err = rows.Scan(&u.ISU, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "scan user")
		}
		users = append(users, u)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "rows error")
	}

	return users, nil
}
