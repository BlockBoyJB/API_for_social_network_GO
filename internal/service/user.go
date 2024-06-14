package service

import (
	"API_for_SN_go/internal/model/pgmodel"
	"API_for_SN_go/internal/repo"
	"API_for_SN_go/internal/repo/pgerrs"
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
)

const (
	userServicePrefixLog = "/service/user"
)

type userService struct {
	userRepo repo.User
}

func newUserService(userRepo repo.User) *userService {
	return &userService{userRepo: userRepo}
}

func (s *userService) UpdateFullName(ctx context.Context, input UserUpdateFullNameInput) error {
	err := s.userRepo.UpdateFullName(ctx, input.Username, input.FirstName, input.LastName)
	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return ErrUserNotFound
		}
		log.Errorf("%s/UpdateFullName error update user full name: %s", userServicePrefixLog, err)
		return ErrCannotUpdateUser
	}
	return nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (pgmodel.User, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return pgmodel.User{}, ErrUserNotFound
		}
		log.Errorf("%s/GetUserByUsername error finding user: %s", userServicePrefixLog, err)
		return pgmodel.User{}, err
	}
	return user, nil
}
