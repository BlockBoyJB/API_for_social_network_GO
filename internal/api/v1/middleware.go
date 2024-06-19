package v1

import (
	"API_for_SN_go/internal/service"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	bearerPrefix = "Bearer "
	usernameCtx  = "username"
)

type AuthMiddleware struct {
	auth service.Auth
}

func (h *AuthMiddleware) AuthHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := parseToken(c.Request())
		if !ok {
			errorResponse(c, http.StatusUnauthorized, ErrInvalidAuthHeader.Error())
			return nil
		}
		username, err := h.auth.ValidateToken(c.Request().Context(), token)
		if err != nil {
			if errors.Is(err, service.ErrCannotParseToken) {
				errorResponse(c, http.StatusUnauthorized, err.Error())
			} else {
				// token is invalid or expired
				errorResponse(c, http.StatusForbidden, err.Error())
			}
			return nil
		}
		c.Set(usernameCtx, username)
		return next(c)
	}
}

// Auth via bearer token in header Authorization
func parseToken(r *http.Request) (string, bool) {
	header := r.Header.Get(echo.HeaderAuthorization)
	if header == "" {
		return "", false
	}
	token := strings.Split(header, bearerPrefix)
	if len(token) == 1 {
		return "", false
	}
	return token[1], true
}

func LoggingMiddleware(h *echo.Echo, output string) {
	cfg := middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339}", "method":"${method}","uri":"${uri}", "status":${status}, "error":"${error}"}` + "\n",
	}
	if output == "stdout" {
		cfg.Output = os.Stdout
	} else {
		file, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			log.Fatal(err)
		}
		cfg.Output = file
	}
	h.Use(middleware.LoggerWithConfig(cfg))
}
