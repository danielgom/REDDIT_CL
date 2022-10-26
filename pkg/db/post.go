package db

import (
	"context"
	"strings"

	"RD-Clone-API/pkg/config"
	"RD-Clone-API/pkg/model"
	"RD-Clone-API/pkg/util/errors"
	"github.com/jackc/pgx/v4/pgxpool"
)

type postRepo struct {
	DB config.DBConn
}

// NewPostRepository creates a new post repository instance.
func NewPostRepository(conn *pgxpool.Pool) PostRepository {
	return &postRepo{DB: conn}
}

// Save creates a new post for a subreddit.
func (p *postRepo) Save(ctx context.Context, post *model.Post) (*model.Post, errors.CommonError) {
	row := p.DB.QueryRow(ctx, `INSERT INTO post("name", "url", "description", "vote_count",
                 "user_id", "created_date", "subreddit_id") VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		post.Name, post.URL, post.Description, post.VoteCount, post.UserID, post.CreatedAt, post.SubredditID)

	saveErr := row.Scan(&post.ID)
	if saveErr != nil {
		return nil, errors.NewInternalServerError("Database error", saveErr)
	}

	return post, nil
}

// FindByID finds a post by ID.
func (p *postRepo) FindByID(ctx context.Context, postID int) (*model.Post, errors.CommonError) {
	var post model.Post

	findErr := p.DB.QueryRow(ctx, `SELECT * FROM post WHERE id = $1`, postID).Scan(&post.ID, &post.Name, &post.URL,
		&post.Description, &post.VoteCount, &post.UserID, &post.CreatedAt, &post.SubredditID)

	if findErr != nil {
		if strings.Contains(findErr.Error(), errNotFound) {
			return nil, errors.NewNotFoundError("user not found")
		}
		return nil, errors.NewInternalServerError("Database error", findErr)
	}

	return &post, nil
}

// FindAllByUser finds all the post from the logged user ordered by date.
func (p *postRepo) FindAllByUser(ctx context.Context, userID int) ([]*model.Post, errors.CommonError) {
	return p.findAllBy(ctx, `SELECT * FROM post WHERE user_id = $1 ORDER BY created_date DESC `, userID)
}

// FindAllBySubreddit finds all the post from a certain subreddit.
func (p *postRepo) FindAllBySubreddit(ctx context.Context, subredditID int) ([]*model.Post, errors.CommonError) {
	return p.findAllBy(ctx, `SELECT * FROM post WHERE subreddit_id = $1`, subredditID)
}

// findAllBy find all the post depending on the query.
func (p *postRepo) findAllBy(ctx context.Context, query string, args ...any) ([]*model.Post, errors.CommonError) {
	var posts []*model.Post

	rows, err := p.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.NewInternalServerError("Database error", err)
	}

	for rows.Next() {
		var post model.Post
		err = rows.Scan(&post.ID, &post.Name, &post.URL,
			&post.Description, &post.VoteCount, &post.UserID, &post.CreatedAt, &post.SubredditID)
		if err != nil {
			return nil, errors.NewInternalServerError("Database error", err)
		}

		posts = append(posts, &post)
	}

	return posts, nil
}
