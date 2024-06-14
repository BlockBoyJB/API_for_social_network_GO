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

const (
	postServicePrefixLog = "/service/post"
)

type postService struct {
	postRepo repo.Post
}

func newPostService(postRepo repo.Post) *postService {
	return &postService{postRepo: postRepo}
}

func (s *postService) CreatePost(ctx context.Context, input PostCreateInput) (string, error) {
	postId := uuid.NewString()
	err := s.postRepo.CreatePost(ctx, pgmodel.Post{
		Username: input.Username,
		PostId:   postId,
		Title:    input.Title,
		Text:     input.Text,
	})
	if err != nil {
		if errors.Is(err, pgerrs.ErrAlreadyExists) {
			return "", ErrPostAlreadyExists // конечно это маловероятно
		}
		if errors.Is(err, pgerrs.ErrForeignKey) {
			return "", ErrUserNotFound // нарушение внешнего ключа возможно только если пользователь не существует
		}
		log.Errorf("%s/CreatePost error create post: %s", postServicePrefixLog, err)
		return "", ErrCannotCreatePost
	}
	return postId, nil
}

func (s *postService) GetPostById(ctx context.Context, postId string) (pgmodel.Post, error) {
	post, err := s.postRepo.GetPostById(ctx, postId)
	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return pgmodel.Post{}, ErrPostNotFound
		}
		log.Errorf("%s/GetPostById error find post by id: %s", postServicePrefixLog, err)
		return pgmodel.Post{}, err
	}
	return post, nil
}
