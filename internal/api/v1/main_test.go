package v1

import (
	"API_for_SN_go/internal/repo"
	"API_for_SN_go/internal/service"
	"API_for_SN_go/pkg/hasher"
	"API_for_SN_go/pkg/postgres"
	"API_for_SN_go/pkg/redis"
	"API_for_SN_go/pkg/validator"
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type APITestSuite struct {
	suite.Suite
	router       *echo.Echo
	pg           *postgres.Postgres
	repositories *repo.Repositories
	redis        *redis.Redis
	services     *service.Services
	m            *migrate.Migrate
}

func (s *APITestSuite) SetupTest() {
	testPGUrl := "postgres://postgres:1234567890@localhost:6000/postgres"
	redisUrl := "127.0.0.1:6379"
	m, err := migrate.New("file://../../../migrations", testPGUrl+"?sslmode=disable")
	if err != nil {
		panic(err)
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
	s.m = m
	pg, err := postgres.NewPG(testPGUrl)
	if err != nil {
		panic(err)
	}
	s.pg = pg

	s.redis = redis.NewRedis(redisUrl)

	s.repositories = repo.NewRepositories(pg)
	d := service.ServicesDependencies{
		Repos:    s.repositories,
		Hasher:   hasher.NewHasher("secret"),
		Redis:    s.redis,
		SignKey:  "secret",
		TokenTTL: time.Hour,
	}
	s.services = service.NewServices(d)

	s.router = echo.New()
	s.router.Validator, err = validator.NewValidator()
	if err != nil {
		panic(err)
	}
	NewRouter(s.router, s.services)
}

// Для тестирования эндпойнтов (кроме /auth группы) нужен авторизованный пользователь (токен)
type apiTestsInfo struct {
	username string
	password string
	token    string
}

func setupApiTests(s *APITestSuite) *apiTestsInfo {
	username, password := "vasek", "1234"
	if err := s.services.Auth.CreateUser(context.Background(), service.UserCreateInput{
		Username:  username,
		FirstName: "Vasya",
		LastName:  "Pupkin",
		Email:     "test",
		Password:  password,
	}); err != nil {
		panic(err)
	}
	token, err := s.services.Auth.CreateToken(context.Background(), service.UserAuthInput{
		Username: username,
		Password: password,
	})
	if err != nil {
		panic(err)
	}
	return &apiTestsInfo{
		username: username,
		password: password,
		token:    token,
	}
}

func tearDownApiTests(s *APITestSuite, setup *apiTestsInfo) {
	if err := s.services.Auth.DeleteUser(context.Background(), service.UserDeleteInput{
		Username: setup.username,
		Password: setup.password,
	}); err != nil {
		panic(err)
	}
	if err := s.redis.Pool.Del(context.Background(), "jwt:"+setup.username).Err(); err != nil {
		panic(err)
	}
}

func (s *APITestSuite) TearDownTest() {
	_ = s.m.Drop()
	s.pg.Close()
	s.redis.Close()
}

func TestAllIntegrations(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
