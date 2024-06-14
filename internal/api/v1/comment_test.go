package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
)

func (s *APITestSuite) Test_commentRouter_create() {
	setup := setupReactionRouterTests(s)
	defer tearDownRouterTests(s, setup)

	testCases := []struct {
		testName   string
		inputBody  string
		expectCode int
	}{
		{
			testName:   "correct test",
			inputBody:  fmt.Sprintf(`{"post_id": "%s", "comment": "good"}`, setup.postId),
			expectCode: 201,
		},
		{
			testName:   "incorrect post id",
			inputBody:  `{"post_id": "0", "comment": "subscribe on my channel"}`,
			expectCode: 400,
		},
	}
	for _, tc := range testCases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/comment/create", bytes.NewBufferString(tc.inputBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+setup.token)
		s.router.ServeHTTP(w, req)
		s.Assert().Equal(tc.expectCode, w.Code)
		if tc.expectCode == 201 {
			var response struct {
				CommentId string `json:"comment_id"`
			}
			_ = json.Unmarshal(w.Body.Bytes(), &response)
			_, err := s.services.Comment.GetCommentById(context.Background(), response.CommentId)
			s.Assert().Equal(nil, err)
		}
	}
}
