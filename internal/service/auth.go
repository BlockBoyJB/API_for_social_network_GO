package service

import (
	"API_for_SN_go/internal/model/pgmodel"
	"API_for_SN_go/internal/repo"
	"API_for_SN_go/internal/repo/pgerrs"
	"API_for_SN_go/pkg/hasher"
	"API_for_SN_go/pkg/redis"
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	defaultKeyPrefix     = "jwt:"
	authServicePrefixLog = "/service/auth"
)

type TokenClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

type authService struct {
	userRepo repo.User
	hasher   hasher.PasswordHasher
	redis    *redis.Redis
	signKey  string
	tokenTTL time.Duration
}

func newAuthService(userRepo repo.User, hasher hasher.PasswordHasher, redis *redis.Redis, singKey string, tokenTTL time.Duration) *authService {
	return &authService{
		userRepo: userRepo,
		hasher:   hasher,
		redis:    redis,
		signKey:  singKey,
		tokenTTL: tokenTTL,
	}
}
func (s *authService) CreateToken(ctx context.Context, input UserAuthInput) (string, error) {
	ok, err := s.verifyPassword(ctx, input.Username, input.Password)
	if !ok || err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Username: input.Username,
	})
	signedToken, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		log.Errorf("%s/CreateToken error sign claims: %s", authServicePrefixLog, err)
		return "", ErrCannotCreateToken
	}
	if err = s.redis.Pool.Set(ctx, defaultKeyPrefix+input.Username, signedToken, s.tokenTTL).Err(); err != nil {
		log.Errorf("%s/CreateToken error save token to redis: %s", authServicePrefixLog, err)
		return "", ErrCannotCreateToken
	}
	return signedToken, nil
}

func (s *authService) ParseToken(ctx context.Context, tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.signKey), nil
	})
	if err != nil {
		log.Errorf("%s/ParseToken error parse token: %s", authServicePrefixLog, err)
		return "", ErrCannotParseToken
	}
	if !token.Valid {
		return "", ErrInvalidToken
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return "", ErrCannotParseToken
	}
	redisToken, err := s.redis.Pool.Get(ctx, defaultKeyPrefix+claims.Username).Result()
	if err != nil || redisToken != tokenString {
		if err != nil {
			log.Errorf("%s/ParseToken error find token from redis: %s", authServicePrefixLog, err)
		}
		return "", ErrExpiredToken
	}
	return claims.Username, nil
}

func (s *authService) CreateUser(ctx context.Context, input UserCreateInput) error {
	err := s.userRepo.CreateUser(ctx, pgmodel.User{
		Username:  input.Username,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  s.hasher.Hash(input.Password),
	})
	if err != nil {
		if errors.Is(err, pgerrs.ErrAlreadyExists) {
			return ErrUserAlreadyExists
		}
		log.Errorf("%s/CreateUser error create user: %s", authServicePrefixLog, err)
		return ErrCannotCreateUser
	}
	return nil
}

func (s *authService) DeleteUser(ctx context.Context, input UserDeleteInput) error {
	ok, err := s.verifyPassword(ctx, input.Username, input.Password)
	if !ok || err != nil {
		return err
	}
	if err = s.userRepo.DeleteUser(ctx, input.Username); err != nil {
		log.Errorf("%s/DeleteUser error delete user: %s", authServicePrefixLog, err)
		return ErrCannotDeleteUser
	}
	if err = s.redis.Pool.Del(ctx, defaultKeyPrefix+input.Username).Err(); err != nil {
		log.Errorf("%s/DeleteUser error delete user token from redis: %s", authServicePrefixLog, err)
		return ErrCannotDeleteUser
	}
	return nil
}

func (s *authService) UpdateUsername(ctx context.Context, input UpdateUsernameInput) error {
	ok, err := s.verifyPassword(ctx, input.Username, input.Password)
	if !ok || err != nil {
		return err
	}
	if err = s.userRepo.UpdateUsername(ctx, input.Username, input.NewUsername); err != nil {
		log.Errorf("%s/UpdateUsername error update username: %s", authServicePrefixLog, err)
		return ErrCannotUpdateUser
	}
	if err = s.redis.Pool.Del(ctx, defaultKeyPrefix+input.Username).Err(); err != nil {
		log.Errorf("%s/UpdateUsername error delete user token from redis: %s", authServicePrefixLog, err)
		return ErrCannotUpdateUser
	}
	return nil
}

func (s *authService) verifyPassword(ctx context.Context, username, password string) (bool, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return false, ErrUserNotFound
		}
		log.Errorf("%s/verifyPassword error verifying user password: %s", authServicePrefixLog, err)
		return false, ErrIncorrectPassword
	}
	if !s.hasher.Verify(password, user.Password) {
		return false, ErrIncorrectPassword
	}
	return true, nil
}
