package service

import (
	"API_for_SN_go/internal/model/pgmodel"
	"API_for_SN_go/internal/repo"
	"API_for_SN_go/pkg/hasher"
	"API_for_SN_go/pkg/redis"
	"context"
	"time"
)

type (
	UserCreateInput struct {
		Username  string
		FirstName string
		LastName  string
		Email     string
		Password  string
	}
	UserAuthInput struct {
		Username string
		Password string
	}
	UserDeleteInput struct {
		Username string
		Password string
	}
	UpdateUsernameInput struct {
		Username    string
		NewUsername string
		Password    string
	}
	UserUpdateFullNameInput struct {
		Username  string
		FirstName string
		LastName  string
	}
	Auth interface {
		CreateToken(ctx context.Context, input UserAuthInput) (string, error)
		ParseToken(ctx context.Context, tokenString string) (string, error)

		CreateUser(ctx context.Context, input UserCreateInput) error
		DeleteUser(ctx context.Context, input UserDeleteInput) error
		UpdateUsername(ctx context.Context, input UpdateUsernameInput) error
	}
	User interface {
		UpdateFullName(ctx context.Context, input UserUpdateFullNameInput) error
		GetUserByUsername(ctx context.Context, username string) (pgmodel.User, error)
	}
)

type (
	PostCreateInput struct {
		Username string
		Title    string
		Text     string
	}
	Post interface {
		CreatePost(ctx context.Context, input PostCreateInput) (string, error)
		GetPostById(ctx context.Context, postId string) (pgmodel.Post, error)
	}
)

type (
	ReactionCreateInput struct {
		PostId   string
		Reaction string
	}
	Reaction interface {
		CreateReaction(ctx context.Context, input ReactionCreateInput) (string, error)
		GetManyReactions(ctx context.Context, postId string) (map[string]string, error)
		GetReactionById(ctx context.Context, reactionId string) (pgmodel.Reaction, error)
		DeleteReaction(ctx context.Context, reactionId string) error
	}
)

type (
	CommentCreateInput struct {
		Username string
		PostId   string
		Comment  string
	}
	CommentUpdateInput struct {
		Username   string
		CommentId  string
		NewComment string
	}
	CommentDeleteInput struct {
		Username  string
		CommentId string
	}
	Comment interface {
		CreateComment(ctx context.Context, input CommentCreateInput) (string, error)
		GetCommentById(ctx context.Context, commentId string) (pgmodel.Comment, error)
		GetManyComments(ctx context.Context, filter, filterParams string) (map[string]string, error)
		UpdateComment(ctx context.Context, input CommentUpdateInput) error
		DeleteComment(ctx context.Context, input CommentDeleteInput) error
	}
)

type (
	Services struct {
		Auth     Auth
		User     User
		Post     Post
		Reaction Reaction
		Comment  Comment
	}
	ServicesDependencies struct {
		Repos    *repo.Repositories
		Hasher   hasher.PasswordHasher
		Redis    *redis.Redis
		SignKey  string
		TokenTTL time.Duration
	}
)

func NewServices(d ServicesDependencies) *Services {
	return &Services{
		Auth:     newAuthService(d.Repos.User, d.Hasher, d.Redis, d.SignKey, d.TokenTTL),
		User:     newUserService(d.Repos.User),
		Post:     newPostService(d.Repos.Post),
		Reaction: newReactionService(d.Repos.Reaction),
		Comment:  newCommentService(d.Repos.Comment),
	}
}
