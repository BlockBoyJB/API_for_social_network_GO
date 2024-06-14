package v1

import (
	"API_for_SN_go/internal/service"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type commentRouter struct {
	commentService service.Comment
}

func newCommentRouter(g *echo.Group, commentService service.Comment) {
	r := &commentRouter{commentService: commentService}
	g.POST("/create", r.create)
	g.PUT("/update", r.updateComment)
	g.DELETE("/delete", r.deleteComment)
	g.GET("", r.getCommentById)
}

type commentCreateInput struct {
	PostId  string `json:"post_id"`
	Comment string `json:"comment"`
}

// @Summary		Create comment
// @Description	Create comment for post
// @Tags			comment
// @Accept			json
// @Produce		json
// @Param			input	body		commentCreateInput	true	"input"
// @Success		201		{object}	map[string]string
// @Failure		400		{object}	echo.HTTPError
// @Failure		500		{object}	echo.HTTPError
// @Security		JWT
// @Router			/api/v1/posts/comment/create [post]
func (r *commentRouter) create(c echo.Context) error {
	var input commentCreateInput

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
	commentId, err := r.commentService.CreateComment(c.Request().Context(), service.CommentCreateInput{
		Username: username,
		PostId:   input.PostId,
		Comment:  input.Comment,
	})
	if err != nil {
		if errors.Is(err, service.ErrPostAlreadyExists) {
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
		CommentId string `json:"comment_id"`
	}
	return c.JSON(http.StatusCreated, response{CommentId: commentId})
}

type commentUpdateInput struct {
	CommentId  string `json:"comment_id"`
	NewComment string `json:"new_comment"`
}

// @Summary		Update comment
// @Description	Update comment for post by commentId
// @Tags			comment
// @Accept			json
// @Produce		json
// @Param			input	body	commentUpdateInput	true	"input"
// @Success		200
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Security		JWT
// @Router			/api/v1/posts/comment/update [put]
func (r *commentRouter) updateComment(c echo.Context) error {
	var input commentUpdateInput

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
	err := r.commentService.UpdateComment(c.Request().Context(), service.CommentUpdateInput{
		Username:   username,
		CommentId:  input.CommentId,
		NewComment: input.NewComment,
	})
	if err != nil {
		if errors.Is(err, service.ErrCommentNotFound) {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}
	return c.NoContent(http.StatusOK)
}

type commentDeleteInput struct {
	CommentId string `json:"comment_id"`
}

// @Summary		Delete comment
// @Description	Delete comment for post
// @Tags			comment
// @Accept			json
// @Produce		json
// @Param			input	body	commentDeleteInput	true	"input"
// @Success		200
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Security		JWT
// @Router			/api/v1/posts/comment/delete [delete]
func (r *commentRouter) deleteComment(c echo.Context) error {
	var input commentDeleteInput

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
	err := r.commentService.DeleteComment(c.Request().Context(), service.CommentDeleteInput{
		Username:  username,
		CommentId: input.CommentId,
	})
	if err != nil {
		if errors.Is(err, service.ErrCommentNotFound) {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}
	return c.NoContent(http.StatusOK)
}

// @Summary		Get comment
// @Description	Get comment for post by commentId
// @Tags			comment
// @Accept			json
// @Produce		json
// @Param			comment_id	query		string	true	"comment id"
// @Success		200			{object}	map[string]string
// @Failure		400			{object}	echo.HTTPError
// @Failure		500			{object}	echo.HTTPError
// @Security		JWT
// @Router			/api/v1/posts/comment [get]
func (r *commentRouter) getCommentById(c echo.Context) error {
	commentId := c.QueryParam("comment_id")
	if len(commentId) == 0 {
		errorResponse(c, http.StatusBadRequest, "invalid request params")
		return nil
	}
	comment, err := r.commentService.GetCommentById(c.Request().Context(), commentId)
	if err != nil {
		if errors.Is(err, service.ErrCommentNotFound) {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}
	type response struct {
		Username  string `json:"username"`
		PostId    string `json:"post_id"`
		CommentId string `json:"comment_id"`
		Comment   string `json:"comment"`
	}
	return c.JSON(http.StatusOK, response{
		Username:  comment.Username,
		PostId:    comment.PostId,
		CommentId: comment.CommentId,
		Comment:   comment.Comment,
	})
}
