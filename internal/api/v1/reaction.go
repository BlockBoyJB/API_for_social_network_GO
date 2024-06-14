package v1

import (
	"API_for_SN_go/internal/service"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type reactionRouter struct {
	reactionService service.Reaction
}

func newReactionRouter(g *echo.Group, reactionService service.Reaction) {
	r := &reactionRouter{reactionService: reactionService}
	g.POST("/create", r.create)
	g.GET("", r.getReactionById)
	g.DELETE("/delete", r.deleteReaction)
}

type reactionCreateInput struct {
	PostId   string `json:"post_id" validate:"required"`
	Reaction string `json:"reaction" validate:"required,reaction"`
}

// @Summary		Create reaction
// @Description	Create reaction for post
// @Tags			reaction
// @Accept			json
// @Produce		json
// @Param			input	body		reactionCreateInput	true	"input"
// @Success		201		{object}	map[string]string
// @Failure		400		{object}	echo.HTTPError
// @Failure		500		{object}	echo.HTTPError
// @Security		JWT
// @Router			/api/v1/posts/reaction/create [post]
func (r *reactionRouter) create(c echo.Context) error {
	var input reactionCreateInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}
	if err := c.Validate(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	reactionId, err := r.reactionService.CreateReaction(c.Request().Context(), service.ReactionCreateInput{
		PostId:   input.PostId,
		Reaction: input.Reaction,
	})
	if err != nil {
		if errors.Is(err, service.ErrReactionAlreadyExists) {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		if errors.Is(err, service.ErrPostNotFound) {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	type response struct {
		ReactionId string `json:"reaction_id"`
	}
	return c.JSON(http.StatusCreated, response{ReactionId: reactionId})
}

// @Summary		Get reaction
// @Description	Get reaction for post by id
// @Tags			reaction
// @Accept			json
// @Produce		json
// @Param			reaction_id	query	string	true	"reaction id"
// @Success		200		{object}	map[string]string
// @Failure		400		{object}	echo.HTTPError
// @Failure		500		{object}	echo.HTTPError
// @Security		JWT
// @Router			/api/v1/posts/reaction [get]
func (r *reactionRouter) getReactionById(c echo.Context) error {
	reactionId := c.QueryParam("reaction_id")
	if len(reactionId) == 0 {
		errorResponse(c, http.StatusBadRequest, "invalid request params")
		return nil
	}
	reaction, err := r.reactionService.GetReactionById(c.Request().Context(), reactionId)
	if err != nil {
		if errors.Is(err, service.ErrReactionNotFound) {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		PostId     string `json:"post_id"`
		ReactionId string `json:"reaction_id"`
		Reaction   string `json:"reaction"`
	}
	return c.JSON(http.StatusOK, response{
		PostId:     reaction.PostId,
		ReactionId: reaction.ReactionId,
		Reaction:   reaction.Reaction,
	})
}

type reactionDeleteInput struct {
	ReactionId string `json:"reaction_id" validate:"required"`
}

// @Summary		Delete reaction
// @Description	Delete reaction for post by id
// @Tags			reaction
// @Accept			json
// @Produce		json
// @Param			input	body		reactionDeleteInput	true	"input"
// @Success		200
// @Failure		400		{object}	echo.HTTPError
// @Failure		500		{object}	echo.HTTPError
// @Security		JWT
// @Router			/api/v1/posts/reaction/delete [delete]
func (r *reactionRouter) deleteReaction(c echo.Context) error {
	var input reactionDeleteInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}
	err := r.reactionService.DeleteReaction(c.Request().Context(), input.ReactionId)
	if err != nil {
		if errors.Is(err, service.ErrReactionNotFound) {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}
	return c.NoContent(http.StatusOK)
}
