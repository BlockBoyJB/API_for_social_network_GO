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

const commentServicePrefixLog = "/service/comment"

type commentService struct {
	commentRepo repo.Comment
}

func newCommentService(commentRepo repo.Comment) *commentService {
	return &commentService{commentRepo: commentRepo}
}

func (s *commentService) CreateComment(ctx context.Context, input CommentCreateInput) (string, error) {
	commentId := uuid.NewString()
	err := s.commentRepo.CreateComment(ctx, pgmodel.Comment{
		Username:  input.Username,
		PostId:    input.PostId,
		CommentId: commentId,
		Comment:   input.Comment,
	})
	if err != nil {
		if errors.Is(err, pgerrs.ErrAlreadyExists) {
			return "", ErrCommentAlreadyExists
		}
		if errors.Is(err, pgerrs.ErrForeignKey) {
			return "", ErrPostNotFound // нарушение внешнего ключа возможно только если пост не существует
		}
		log.Errorf("%s/CreateComment error create comment: %s", commentServicePrefixLog, err)
		return "", ErrCannotCreateComment
	}
	return commentId, nil
}

func (s *commentService) GetCommentById(ctx context.Context, commentId string) (pgmodel.Comment, error) {
	comment, err := s.commentRepo.GetCommentById(ctx, commentId)
	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return pgmodel.Comment{}, ErrCommentNotFound
		}
		return pgmodel.Comment{}, err
	}
	return comment, nil
}

func (s *commentService) GetManyComments(ctx context.Context, filter, filterParams string) (map[string]string, error) {
	comments, err := s.commentRepo.GetManyComments(ctx, filter, filterParams)
	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return nil, ErrCommentNotFound
		}
		log.Errorf("%s/GetManyComments error finding comments: %s", commentServicePrefixLog, err)
		return nil, err
	}
	res := make(map[string]string)
	for _, comment := range comments {
		res[comment.CommentId] = comment.Comment
	}
	return res, nil
}

func (s *commentService) UpdateComment(ctx context.Context, input CommentUpdateInput) error {
	err := s.commentRepo.UpdateComment(ctx, input.Username, input.CommentId, input.NewComment)
	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return ErrCommentNotFound
		}
		return ErrCannotCreateComment
	}
	return nil
}

func (s *commentService) DeleteComment(ctx context.Context, input CommentDeleteInput) error {
	err := s.commentRepo.DeleteComment(ctx, input.Username, input.CommentId)
	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return ErrCommentNotFound
		}
		return ErrCannotDeleteComment
	}
	return nil
}
