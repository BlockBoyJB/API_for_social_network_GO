package repo

import (
	"API_for_SN_go/internal/model/pgmodel"
	"API_for_SN_go/internal/repo/pgdb"
	"API_for_SN_go/pkg/postgres"
	"context"
)

type User interface {
	CreateUser(ctx context.Context, u pgmodel.User) error
	GetUserByUsername(ctx context.Context, username string) (pgmodel.User, error)
	UpdateUsername(ctx context.Context, username, newUsername string) error
	UpdateFullName(ctx context.Context, username, firstName, lastName string) error
	DeleteUser(ctx context.Context, username string) error
}

type Post interface {
	CreatePost(ctx context.Context, p pgmodel.Post) error
	GetPostById(ctx context.Context, postId string) (pgmodel.Post, error)
}

type Reaction interface {
	CreateReaction(ctx context.Context, rn pgmodel.Reaction) error
	GetReactionById(ctx context.Context, reactionId string) (pgmodel.Reaction, error)
	GetManyReactions(ctx context.Context, postId string) ([]pgmodel.Reaction, error)
	DeleteReaction(ctx context.Context, reactionId string) error
}

type Comment interface {
	CreateComment(ctx context.Context, c pgmodel.Comment) error
	GetCommentById(ctx context.Context, commentId string) (pgmodel.Comment, error)
	GetManyComments(ctx context.Context, filter, filterParams string) ([]pgmodel.Comment, error)
	UpdateComment(ctx context.Context, username, commentId, newComment string) error
	DeleteComment(ctx context.Context, username, commentId string) error
}

type Repositories struct {
	User
	Post
	Reaction
	Comment
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User:     pgdb.NewUserRepo(pg),
		Post:     pgdb.NewPostRepo(pg),
		Reaction: pgdb.NewReactionRepo(pg),
		Comment:  pgdb.NewCommentRepo(pg),
	}
}
