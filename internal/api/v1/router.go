package v1

import (
	_ "API_for_SN_go/docs"
	"API_for_SN_go/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewRouter(h *echo.Echo, services *service.Services) {
	h.Use(middleware.Recover())
	h.GET("/ping", ping)
	h.GET("/swagger/*", echoSwagger.WrapHandler)

	newAuthRouter(h.Group("/auth"), services.Auth)
	authMiddleware := &AuthMiddleware{auth: services.Auth}
	v1 := h.Group("/api/v1", authMiddleware.AuthHandler)

	newUserRouter(v1.Group("/user"), services.User, services.Comment)
	newPostRouter(v1.Group("/posts/post"), services.Post, services.Reaction, services.Comment)
	newReactionRouter(v1.Group("/posts/reaction"), services.Reaction)
	newCommentRouter(v1.Group("/posts/comment"), services.Comment)
}

func ping(c echo.Context) error {
	return c.NoContent(200)
}
