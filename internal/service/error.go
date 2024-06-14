package service

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrCannotCreateUser  = errors.New("cannot create user")
	ErrCannotDeleteUser  = errors.New("cannot delete user")
	ErrCannotUpdateUser  = errors.New("cannot update user info")
	ErrUserNotFound      = errors.New("user not found")
	ErrIncorrectPassword = errors.New("incorrect user password")

	ErrCannotCreateToken = errors.New("cannot create token")
	ErrInvalidToken      = errors.New("invalid token")
	ErrExpiredToken      = errors.New("expired token")
	ErrCannotParseToken  = errors.New("cannot parse token")

	ErrCannotCreatePost  = errors.New("cannot create post")
	ErrPostAlreadyExists = errors.New("post already exists")
	ErrPostNotFound      = errors.New("post not found")

	ErrReactionAlreadyExists = errors.New("reaction already exists")
	ErrReactionNotFound      = errors.New("reaction not found")
	ErrCannotCreateReaction  = errors.New("cannot create reaction")

	ErrCommentAlreadyExists = errors.New("comment already exists")
	ErrCannotCreateComment  = errors.New("cannot create comment")
	ErrCommentNotFound      = errors.New("comment not found")
	ErrCannotDeleteComment  = errors.New("cannot delete comment")
)
