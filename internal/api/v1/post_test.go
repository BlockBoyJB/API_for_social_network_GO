package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
)

func (s *APITestSuite) Test_postRouterCreate() {
	setup := setupApiTests(s)
	defer tearDownApiTests(s, setup)

	testCases := []struct {
		testName   string
		inputBody  string
		expectCode int
	}{
		{
			testName:   "correct test",
			inputBody:  `{"title": "test_title", "text": "test_text"}`,
			expectCode: 201,
		},
	}
	for _, tc := range testCases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/post/create", bytes.NewBufferString(tc.inputBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+setup.token)
		s.router.ServeHTTP(w, req)
		s.Assert().Equal(tc.expectCode, w.Code)

		if tc.expectCode == 201 {
			var response struct {
				PostId string `json:"post_id"`
			}
			_ = json.Unmarshal(w.Body.Bytes(), &response)
			_, err := s.services.Post.GetPostById(context.Background(), response.PostId)
			s.Assert().Equal(nil, err)
		}
	}
}
