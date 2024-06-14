package v1

import (
	"API_for_SN_go/internal/service"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type postRouter struct {
	postService     service.Post
	reactionService service.Reaction
	commentService  service.Comment
}

func newPostRouter(g *echo.Group, postService service.Post, reactionService service.Reaction, commentService service.Comment) {
	r := &postRouter{
		postService:     postService,
		reactionService: reactionService,
		commentService:  commentService,
	}
	g.POST("/create", r.create)
	g.GET("", r.getById)
	g.GET("/comments", r.getPostComments)
}

type postCreateInput struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// @Summary		Create post
// @Description	Create post
// @Tags			post
// @Accept			json
// @Produce		json
// @Param			input	body		postCreateInput	true	"input"
// @Success		201		{object}	map[string]string
// @Failure		400		{object}	echo.HTTPError
// @Failure		500		{object}	echo.HTTPError
// @Security		JWT
// @Router			/api/v1/posts/post/create [post]
func (r *postRouter) create(c echo.Context) error {
	var input postCreateInput

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
	postId, err := r.postService.CreatePost(c.Request().Context(), service.PostCreateInput{
		Username: username,
		Title:    input.Title,
		Text:     input.Text,
	})
	if err != nil {
		if errors.Is(err, service.ErrPostAlreadyExists) {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		if errors.Is(err, service.ErrUserNotFound) {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return err
	}

	type response struct {
		PostId string `json:"post_id"`
	}
	return c.JSON(http.StatusCreated, response{PostId: postId})
}

// @Summary		Get post
// @Description	Get post by id
// @Tags			post
// @Accept			json
// @Produce		json
// @Param			post_id	query		string	true	"post id"
// @Success		200		{object}	map[string]string
// @Failure		400		{object}	echo.HTTPError
// @Failure		500		{object}	echo.HTTPError
// @Security		JWT
// @Router			/api/v1/posts/post [get]
func (r *postRouter) getById(c echo.Context) error {
	postId := c.QueryParam("post_id")
	if len(postId) == 0 {
		errorResponse(c, http.StatusBadRequest, "invalid request params")
		return nil
	}
	post, err := r.postService.GetPostById(c.Request().Context(), postId)
	if err != nil {
		if errors.Is(err, service.ErrPostNotFound) {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}
	reactions, err := r.reactionService.GetManyReactions(c.Request().Context(), postId)
	if err != nil {
		if !errors.Is(err, service.ErrReactionNotFound) {
			errorResponse(c, http.StatusInternalServerError, "internal server error")
			return err
		}

	}
	type response struct {
		Username  string            `json:"username"`
		PostId    string            `json:"post_id"`
		Title     string            `json:"title"`
		Text      string            `json:"text"`
		Reactions map[string]string `json:"reactions"`
	}
	return c.JSON(http.StatusOK, response{
		Username:  post.Username,
		PostId:    post.PostId,
		Title:     post.Title,
		Text:      post.Text,
		Reactions: reactions,
	})
}

// @Summary		Get post comments
// @Description	Get all post comments by post id
// @Tags			post
// @Accept			json
// @Produce		json
// @Param			post_id	query		string	true	"post id"
// @Success		200		{object}	map[string]string
// @Failure		400		{object}	echo.HTTPError
// @Failure		500		{object}	echo.HTTPError
// @Security		JWT
// @Router			/api/v1/posts/post/comments [get]
func (r *postRouter) getPostComments(c echo.Context) error {
	postId := c.QueryParam("post_id")
	if len(postId) == 0 {
		errorResponse(c, http.StatusBadRequest, "invalid request params")
		return nil
	}
	comments, err := r.commentService.GetManyComments(c.Request().Context(), "post_id", postId)
	if err != nil {
		if errors.Is(err, service.ErrCommentNotFound) {
			errorResponse(c, http.StatusBadRequest, "comments not found")
			return err
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}
	type response struct {
		PostId   string            `json:"post_id"`
		Comments map[string]string `json:"comments"`
	}
	return c.JSON(http.StatusOK, response{
		PostId:   postId,
		Comments: comments,
	})
}
