package v1

import (
	"API_for_SN_go/internal/service"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type userRouter struct {
	userService    service.User
	commentService service.Comment
}

func newUserRouter(g *echo.Group, userService service.User, commentService service.Comment) {
	r := &userRouter{
		userService:    userService,
		commentService: commentService,
	}
	g.GET("/comments", r.getUserComments)
	g.GET("", r.getUser)
	g.PUT("/update/full-name", r.updateFullName)
}

type userUpdateFullNameInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// @Summary		Update user full name
// @Description	Update user full name
// @Tags			user
// @Accept			json
// @Produce		json
// @Param			input	body	userUpdateFullNameInput	true	"input"
// @Success		200
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Security		JWT
// @Router			/api/v1/user/update/full-name [put]
func (r *userRouter) updateFullName(c echo.Context) error {
	var input userUpdateFullNameInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}
	userCtx := c.Get(usernameCtx)
	username, ok := userCtx.(string)
	if !ok {
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return nil
	}
	err := r.userService.UpdateFullName(c.Request().Context(), service.UserUpdateFullNameInput{
		Username:  username,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}
	return c.NoContent(http.StatusOK)
}

// @Summary		Get user
// @Description	Get user by username
// @Tags			user
// @Accept			json
// @Produce		json
// @Param			username	query		string	true	"username"
// @Success		200			{object}	map[string]string
// @Failure		400			{object}	echo.HTTPError
// @Failure		500			{object}	echo.HTTPError
// @Security		JWT
// @Router			/api/v1/user [get]
func (r *userRouter) getUser(c echo.Context) error {
	username := c.QueryParam("username")
	if len(username) == 0 {
		errorResponse(c, http.StatusBadRequest, "invalid request params")
		return nil
	}
	user, err := r.userService.GetUserByUsername(c.Request().Context(), username)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}
	type response struct {
		Username  string `json:"username"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
	return c.JSON(http.StatusOK, response{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	})
}

// @Summary		Get user comments
// @Description	Get all user comments
// @Tags			user
// @Accept			json
// @Produce		json
// @Param			username	query		string	false	"username"
// @Success		200			{object}	map[string]string
// @Failure		400			{object}	echo.HTTPError
// @Failure		500			{object}	echo.HTTPError
// @Security		JWT
// @Router			/api/v1/user/comments [get]
func (r *userRouter) getUserComments(c echo.Context) error {
	var u string
	username := c.QueryParam("username")
	if len(username) != 0 {
		u = username
	} else {
		userCtx := c.Get(usernameCtx)
		username, ok := userCtx.(string)
		if !ok {
			errorResponse(c, http.StatusInternalServerError, "internal server error")
			return nil
		}
		u = username
	}

	comments, err := r.commentService.GetManyComments(c.Request().Context(), "username", u)
	if err != nil {
		if errors.Is(err, service.ErrCommentNotFound) {
			errorResponse(c, http.StatusBadRequest, "comments not found")
			return err
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}
	type response struct {
		Comments map[string]string `json:"comments"`
	}
	return c.JSON(http.StatusOK, response{Comments: comments})
}
