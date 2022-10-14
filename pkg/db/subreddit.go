package db

import (
	"context"

	"RD-Clone-API/pkg/config"
	"RD-Clone-API/pkg/model"
	"RD-Clone-API/pkg/util/errors"
	"github.com/jackc/pgx/v4/pgxpool"
)

type subredditRepo struct {
	DB config.DBConn
}

// NewSubredditRepository creates a new subreddit repository instance.
func NewSubredditRepository(conn *pgxpool.Pool) SubredditRepository {
	return &subredditRepo{DB: conn}
}

// Save saves a subreddit.
func (s *subredditRepo) Save(ctx context.Context, subreddit *model.Subreddit) (*model.Subreddit, errors.CommonError) {
	row := s.DB.QueryRow(ctx, `INSERT INTO subreddit("name", "description", "created_at", "user_id") VALUES ($1, $2, $3, $4)
		RETURNING id`,
		subreddit.Name, subreddit.Description, subreddit.CreatedAt, subreddit.UserID)

	err := row.Scan(&subreddit.ID)

	if err != nil {
		return nil, errors.NewInternalServerError("Database error", err)
	}

	return subreddit, nil
}

// FindByID looks for a subreddit by id.
func (s *subredditRepo) FindByID(ctx context.Context, id int) (*model.Subreddit, errors.CommonError) {
	var subreddit model.Subreddit
	err := s.DB.QueryRow(ctx, `SELECT * FROM subreddit WHERE id = $1`, id).Scan(&subreddit.ID, &subreddit.Name,
		&subreddit.Description, &subreddit.CreatedAt, &subreddit.UserID)

	if err != nil {
		return nil, errors.NewInternalServerError("Database error", err)
	}

	return &subreddit, nil
}

// FindAll lists all subreddits.
func (s *subredditRepo) FindAll(ctx context.Context) ([]*model.Subreddit, errors.CommonError) {
	var subreddits []*model.Subreddit
	rows, err := s.DB.Query(ctx, `SELECT * FROM subreddit`)
	if err != nil {
		return nil, errors.NewInternalServerError("Database error", err)
	}
	defer rows.Close()

	for rows.Next() {
		var subreddit model.Subreddit
		err = rows.Scan(&subreddit.ID, &subreddit.Name, &subreddit.Description,
			&subreddit.CreatedAt, &subreddit.UserID)
		if err != nil {
			return nil, errors.NewInternalServerError("Database error", err)
		}
		subreddits = append(subreddits, &subreddit)
	}

	return subreddits, nil
}
