package pgdb

import (
	"API_for_SN_go/internal/model/pgmodel"
	"API_for_SN_go/internal/repo/pgerrs"
	"API_for_SN_go/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	log "github.com/sirupsen/logrus"
)

const commentPrefixLog = "/pgdb/comment"

type CommentRepo struct {
	*postgres.Postgres
}

func NewCommentRepo(pg *postgres.Postgres) *CommentRepo {
	return &CommentRepo{pg}
}

func (r *CommentRepo) CreateComment(ctx context.Context, c pgmodel.Comment) error {
	sql, args, _ := r.Builder.
		Insert("comment").
		Columns("username", "post_id", "comment_id", "comment").
		Values(c.Username, c.PostId, c.CommentId, c.Comment).
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
		log.Errorf("%s/CreateComment error exec stmt: %s", commentPrefixLog, err)
		return err
	}
	return nil
}

func (r *CommentRepo) GetCommentById(ctx context.Context, commentId string) (pgmodel.Comment, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("comment").
		Where("comment_id = ?", commentId).
		ToSql()

	var comment pgmodel.Comment

	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&comment.Id,
		&comment.Username,
		&comment.PostId,
		&comment.CommentId,
		&comment.Comment,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgmodel.Comment{}, pgerrs.ErrNotFound
		}
		log.Errorf("%s/GetCommentById error finding comment: %s", commentPrefixLog, err)
		return pgmodel.Comment{}, err
	}
	return comment, nil
}

func (r *CommentRepo) GetManyComments(ctx context.Context, filter, filterParams string) ([]pgmodel.Comment, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("comment").
		Where(fmt.Sprintf("%s = ?", filter), filterParams).
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgerrs.ErrNotFound
		}
		log.Errorf("%s/GetManyComments error finding reactions by %s: %s", commentPrefixLog, filter, err)
		return nil, err
	}
	var comments []pgmodel.Comment
	for rows.Next() {
		var comment pgmodel.Comment
		err = rows.Scan(&comment.Id, &comment.Username, &comment.PostId, &comment.CommentId, &comment.Comment)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, pgerrs.ErrNotFound
			}
			log.Errorf("%s/GetManyComments error finding comments by %s: %s", commentPrefixLog, filter, err)
			continue
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *CommentRepo) UpdateComment(ctx context.Context, username, commentId, newComment string) error {
	sql, args, _ := r.Builder.
		Update("comment").
		Set("comment", newComment).
		Where("username = ? AND comment_id = ?", username, commentId).
		ToSql()
	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgerrs.ErrNotFound
		}
		log.Errorf("%s/UpdateComment error exec stmt: %s", commentPrefixLog, err)
		return err
	}
	return nil
}

func (r *CommentRepo) DeleteComment(ctx context.Context, username, commentId string) error {
	sql, args, _ := r.Builder.
		Delete("comment").
		Where("username = ? AND comment_id = ?", username, commentId).
		ToSql()
	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgerrs.ErrNotFound
		}
		log.Errorf("%s/DeleteComment error exec stmt: %s", commentPrefixLog, err)
		return err
	}
	return nil
}
