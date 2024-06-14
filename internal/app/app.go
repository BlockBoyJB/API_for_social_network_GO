package app

import (
	"API_for_SN_go/config"
	v1 "API_for_SN_go/internal/api/v1"
	"API_for_SN_go/internal/repo"
	"API_for_SN_go/internal/service"
	"API_for_SN_go/pkg/hasher"
	"API_for_SN_go/pkg/httpserver"
	"API_for_SN_go/pkg/postgres"
	"API_for_SN_go/pkg/redis"
	"API_for_SN_go/pkg/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"os"
	"os/signal"
	"syscall"
)

//	@title			Api for social network
//	@version		1.0
//	@description	Api for social networking. Include posts, reactions and comments

//	@host		localhost:8080
//	@BasePath	/

//	@securityDefinitions.apikey	JWT
//	@in							header
//	@name						Authorization
//	@description				JWT token

func Run() {
	// config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	// set up json logger
	setLogger(cfg.Log.Level, cfg.Log.Output)

	// postgresql database
	pg, err := postgres.NewPG(cfg.PG.Url, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Fatalf("Initializing postgres error: %s", err)
	}
	defer pg.Close()

	// pg repositories
	repos := repo.NewRepositories(pg)

	// redis database for redis-jwt strategy
	rdb := redis.NewRedis(cfg.Redis.Url, redis.MaxPoolSize(cfg.Redis.MaxPoolSize))
	defer rdb.Close()

	dependencies := service.ServicesDependencies{
		Repos:    repos,
		Hasher:   hasher.NewHasher(cfg.Hasher.Salt),
		Redis:    rdb,
		SignKey:  cfg.JWT.SignKey,
		TokenTTL: cfg.JWT.TokenTTL,
	}
	services := service.NewServices(dependencies)

	// validator for incoming requests
	v, err := validator.NewValidator()
	if err != nil {
		log.Fatalf("Initializing handler validator error: %s", err)
	}

	// main handler
	handler := echo.New()
	handler.Validator = v
	v1.LoggingMiddleware(handler, cfg.Log.Output)
	v1.NewRouter(handler, services)

	// http server
	httpServer := httpserver.NewServer(handler, httpserver.Port(cfg.HTTP.Port))
	log.Infof("App started! Listening port %s", cfg.HTTP.Port)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app run, signal " + s.String())

	case err = <-httpServer.Notify():
		log.Errorf("/app/run server notify error: %s", err)
	}

	// graceful shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Errorf("/app/run server shutdown error: %s", err)
	}
	log.Infof("App shutdown with exit code 0")
}

// loading environment params from .env
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("load env file error: %s", err)
	}
}
