package service

import (
	"API_for_SN_go/internal/model/pgmodel"
	"API_for_SN_go/internal/repo"
	"API_for_SN_go/internal/repo/pgerrs"
	"context"
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const reactionServicePrefixLog = "/service/reaction"

type reactionService struct {
	reactionRepo repo.Reaction
}

func newReactionService(reactionRepo repo.Reaction) *reactionService {
	return &reactionService{reactionRepo: reactionRepo}
}

func (s *reactionService) CreateReaction(ctx context.Context, input ReactionCreateInput) (string, error) {
	reactionId := uuid.NewString()
	err := s.reactionRepo.CreateReaction(ctx, pgmodel.Reaction{
		PostId:     input.PostId,
		ReactionId: reactionId,
		Reaction:   input.Reaction,
	})
	if err != nil {
		if errors.Is(err, pgerrs.ErrAlreadyExists) {
			return "", ErrReactionAlreadyExists
		}
		if errors.Is(err, pgerrs.ErrForeignKey) {
			return "", ErrPostNotFound // нарушение внешнего ключа возможно только если пост не существует
		}
		log.Errorf("%s/CreateReaction error create reaction: %s", reactionServicePrefixLog, err)
		return "", ErrCannotCreateReaction
	}
	return reactionId, nil
}

func (s *reactionService) GetManyReactions(ctx context.Context, postId string) (map[string]string, error) {
	reactions, err := s.reactionRepo.GetManyReactions(ctx, postId)
	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return nil, ErrReactionNotFound
		}
		log.Errorf("%s/GetManyReactions error find reactions for post: %s", reactionServicePrefixLog, err)
		return nil, ErrReactionNotFound
	}
	res := make(map[string]string)
	for _, reaction := range reactions {
		res[reaction.ReactionId] = reaction.Reaction
	}
	return res, nil
}

func (s *reactionService) GetReactionById(ctx context.Context, reactionId string) (pgmodel.Reaction, error) {
	reaction, err := s.reactionRepo.GetReactionById(ctx, reactionId)
	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return pgmodel.Reaction{}, ErrReactionNotFound
		}
		log.Errorf("%s/GetReactionById error find reaction by id: %s", reactionServicePrefixLog, err)
		return pgmodel.Reaction{}, err
	}
	return reaction, nil
}

func (s *reactionService) DeleteReaction(ctx context.Context, reactionId string) error {
	err := s.reactionRepo.DeleteReaction(ctx, reactionId)
	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return ErrReactionNotFound
		}
		log.Errorf("%s/DeleteReaction error delete reaction: %s", reactionServicePrefixLog, err)
		return err
	}
	return nil
}
