package v1

import (
	"API_for_SN_go/internal/service"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type authRouter struct {
	authService service.Auth
}

func newAuthRouter(g *echo.Group, authService service.Auth) {
	r := &authRouter{authService: authService}
	g.POST("/sign-up", r.signUp)
	g.POST("/sign-in", r.signIn)
	g.DELETE("/user/delete", r.deleteUser)
	g.PUT("/user/update/username", r.updateUsername)
}

type signUpInput struct {
	Username  string `json:"username" validate:"required,username"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

// @Summary		Sign up
// @Description	Sign up
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			input	body	signUpInput	true	"input"
// @Success		201
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/auth/sign-up [post]
func (r *authRouter) signUp(c echo.Context) error {
	var input signUpInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}
	if err := c.Validate(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}
	if err := r.authService.CreateUser(c.Request().Context(), service.UserCreateInput{
		Username:  input.Username,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  input.Password,
	}); err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}
	return c.NoContent(http.StatusCreated)
}

type signInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Summary		Sign in
// @Description	Sign in
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			input	body		signInInput	true	"input"
// @Success		200		{object} map[string]string
// @Failure		400		{object}	echo.HTTPError
// @Failure		500		{object}	echo.HTTPError
// @Router			/auth/sign-in [post]
func (r *authRouter) signIn(c echo.Context) error {
	var input signInInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}
	token, err := r.authService.CreateToken(c.Request().Context(), service.UserAuthInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrCannotCreateToken) {
			errorResponse(c, http.StatusInternalServerError, "internal server error")
			return err
		}
		errorResponse(c, http.StatusForbidden, err.Error())
		return err
	}

	type response struct {
		Token string `json:"token"`
	}
	return c.JSON(http.StatusOK, response{Token: token})
}

// @Summary		Delete user
// @Description	Delete  user
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			input	body	signInInput	true	"input"
// @Success		200
// @Failure		400	{object}	echo.HTTPError
// @Failure		403	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/auth/user/delete [delete]
func (r *authRouter) deleteUser(c echo.Context) error {
	var input signInInput // same fields

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := r.authService.DeleteUser(c.Request().Context(), service.UserDeleteInput{
		Username: input.Username,
		Password: input.Password,
	}); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		if errors.Is(err, service.ErrIncorrectPassword) {
			errorResponse(c, http.StatusForbidden, err.Error())
			return err
		}
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}
	return c.NoContent(http.StatusOK)
}

type updateUsernameInput struct {
	Username    string `json:"username"`
	NewUsername string `json:"new_username" validate:"username"`
	Password    string `json:"password"`
}

// @Summary		Update username
// @Description	Update user username
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			input	body	updateUsernameInput	true	"input"
// @Success		200
// @Response		400	{object}	echo.HTTPError
// @Failure		403	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/user/update/username [put]
func (r *authRouter) updateUsername(c echo.Context) error {
	var input updateUsernameInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}
	if err := c.Validate(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}
	if err := r.authService.UpdateUsername(c.Request().Context(), service.UpdateUsernameInput{
		Username:    input.Username,
		NewUsername: input.NewUsername,
		Password:    input.Password,
	}); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		if errors.Is(err, service.ErrIncorrectPassword) {
			errorResponse(c, http.StatusForbidden, err.Error())
			return err
		}
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}
	return c.NoContent(http.StatusOK)
}
