package pgdb

import (
	"API_for_SN_go/internal/model/pgmodel"
	"API_for_SN_go/internal/repo/pgerrs"
	"API_for_SN_go/pkg/postgres"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	log "github.com/sirupsen/logrus"
)

const postPrefixLog = "/pgdb/post"

type PostRepo struct {
	*postgres.Postgres
}

func NewPostRepo(pg *postgres.Postgres) *PostRepo {
	return &PostRepo{pg}
}

func (r *PostRepo) CreatePost(ctx context.Context, p pgmodel.Post) error {
	sql, args, _ := r.Builder.
		Insert("post").
		Columns("username", "post_id", "title", "text").
		Values(p.Username, p.PostId, p.Title, p.Text).
		ToSql()
	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return pgerrs.ErrAlreadyExists
			}
			if pgErr.Code == "23503" {
				return pgerrs.ErrForeignKey
			}
		}
		log.Errorf("%s/CreatePost error exec stmt: %s", postPrefixLog, err)
		return err
	}
	return nil
}

func (r *PostRepo) GetPostById(ctx context.Context, postId string) (pgmodel.Post, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("post").
		Where("post_id = ?", postId).
		ToSql()

	var post pgmodel.Post

	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&post.Id,
		&post.Username,
		&post.PostId,
		&post.Title,
		&post.Text,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgmodel.Post{}, pgerrs.ErrNotFound
		}
		log.Errorf("%s/GetPostById error finding post: %s", postPrefixLog, err)
		return pgmodel.Post{}, err
	}
	return post, nil
}
