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

const reactionPrefixLog = "/pgdb/reaction"

type ReactionRepo struct {
	*postgres.Postgres
}

func NewReactionRepo(pg *postgres.Postgres) *ReactionRepo {
	return &ReactionRepo{pg}
}

func (r *ReactionRepo) CreateReaction(ctx context.Context, rn pgmodel.Reaction) error {
	sql, args, _ := r.Builder.
		Insert("reaction").
		Columns("post_id", "reaction_id", "reaction").
		Values(rn.PostId, rn.ReactionId, rn.Reaction).
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
		log.Errorf("%s/CreateReaction error exec stmt: %s", reactionPrefixLog, err)
		return err
	}
	return nil
}

func (r *ReactionRepo) GetReactionById(ctx context.Context, reactionId string) (pgmodel.Reaction, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("reaction").
		Where("reaction_id = ?", reactionId).
		ToSql()

	var reaction pgmodel.Reaction
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&reaction.Id,
		&reaction.PostId,
		&reaction.ReactionId,
		&reaction.Reaction,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgmodel.Reaction{}, pgerrs.ErrNotFound
		}
		log.Errorf("%s/GetReactionById error finding reaction: %s", reactionPrefixLog, err)
		return pgmodel.Reaction{}, err
	}
	return reaction, nil
}

func (r *ReactionRepo) GetManyReactions(ctx context.Context, postId string) ([]pgmodel.Reaction, error) {
	sql, args, _ := r.Builder.
		Select("reaction_id", "reaction").
		From("reaction").
		Where("post_id = ?", postId).
		ToSql()
	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgerrs.ErrNotFound
		}
		log.Errorf("%s/GetManyReactions error finding reactions by post id: %s", reactionPrefixLog, err)
		return nil, err
	}
	var reactions []pgmodel.Reaction
	for rows.Next() {
		var reaction pgmodel.Reaction
		err = rows.Scan(&reaction.ReactionId, &reaction.Reaction)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, pgerrs.ErrNotFound
			}
			log.Errorf("%s/GetManyReactions error finding reactions by post id: %s", reactionPrefixLog, err)
			continue
		}
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}

func (r *ReactionRepo) DeleteReaction(ctx context.Context, reactionId string) error {
	sql, args, _ := r.Builder.
		Delete("reaction").
		Where("reaction_id = ?", reactionId).
		ToSql()
	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgerrs.ErrNotFound
		}
		log.Errorf("%s/DeleteReaction error exec stmt: %s", reactionPrefixLog, err)
		return err
	}
	return nil
}
