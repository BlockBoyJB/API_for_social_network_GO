package v1

import (
	"API_for_SN_go/internal/mocks/servicemocks"
	"API_for_SN_go/internal/service"
	"API_for_SN_go/pkg/validator"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReactionRouter_create(t *testing.T) {
	type args struct {
		ctx   context.Context
		input service.ReactionCreateInput
	}
	type MockBehaviour func(m *servicemocks.MockReaction, args args)

	testCases := []struct {
		testName      string
		args          args
		inputBody     string
		mockBehaviour MockBehaviour
		expectCode    int
		expectBody    string
	}{
		{
			testName: "correct test",
			args: args{
				ctx: context.Background(),
				input: service.ReactionCreateInput{
					PostId:   "1000",
					Reaction: "like",
				},
			},
			inputBody: `{"post_id": "1000", "reaction": "like"}`,
			mockBehaviour: func(m *servicemocks.MockReaction, args args) {
				m.EXPECT().CreateReaction(args.ctx, args.input).Return("1234567890", nil)
			},
			expectCode: 201,
			expectBody: `{"reaction_id":"1234567890"}` + "\n",
		},
		{
			testName:      "incorrect reaction input",
			inputBody:     `{"post_id": "1000", "reaction": "321boom@"}`,
			mockBehaviour: func(m *servicemocks.MockReaction, args args) {},
			expectCode:    400,
			expectBody:    `{"message":"field reaction can only consist of lower Latin characters. Min length is 1, max: 16"}` + "\n",
		},
		{
			testName:      "invalid post id",
			inputBody:     `{"reaction": "boom"}`,
			mockBehaviour: func(m *servicemocks.MockReaction, args args) {},
			expectCode:    400,
			expectBody:    `{"message":"field PostId is invalid"}` + "\n",
		},
		{
			testName:      "too long reaction input",
			inputBody:     `{"post_id": "1000", "reaction": "loooooooooooooooooooooooooongboom"}`,
			mockBehaviour: func(m *servicemocks.MockReaction, args args) {},
			expectCode:    400,
			expectBody:    `{"message":"field reaction can only consist of lower Latin characters. Min length is 1, max: 16"}` + "\n",
		},
		{
			testName:      "too short reaction input (without reaction)",
			inputBody:     `{"post_id": "1000", "reaction": ""}`,
			mockBehaviour: func(m *servicemocks.MockReaction, args args) {},
			expectCode:    400,
			expectBody:    `{"message":"field Reaction is invalid"}` + "\n",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			reaction := servicemocks.NewMockReaction(ctrl)
			tc.mockBehaviour(reaction, tc.args)
			services := &service.Services{Reaction: reaction}

			e := echo.New()
			e.Validator, _ = validator.NewValidator()
			newReactionRouter(e.Group("/api/v1/posts/reaction"), services.Reaction)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/reaction/create", bytes.NewBufferString(tc.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			e.ServeHTTP(w, req)

			assert.Equal(t, tc.expectCode, w.Code)
			assert.Equal(t, tc.expectBody, w.Body.String())
		})
	}
}

type reactionTestInfo struct {
	*apiTestsInfo
	postId string
}

func setupReactionRouterTests(s *APITestSuite) *reactionTestInfo {
	apiSetup := setupApiTests(s)
	postId, err := s.services.Post.CreatePost(context.Background(), service.PostCreateInput{
		Username: apiSetup.username,
		Title:    "testtitle",
		Text:     "test",
	})
	if err != nil {
	}
	return &reactionTestInfo{
		apiTestsInfo: apiSetup,
		postId:       postId,
	}
}

func tearDownRouterTests(s *APITestSuite, setup *reactionTestInfo) {
	tearDownApiTests(s, setup.apiTestsInfo)
	// посты удалятся после удаления пользователя
}

func (s *APITestSuite) Test_reactionRouter_create() {
	setup := setupReactionRouterTests(s)
	defer tearDownRouterTests(s, setup)

	testCases := []struct {
		testName   string
		inputBody  string
		expectCode int
	}{
		{
			testName:   "correct test",
			inputBody:  fmt.Sprintf(`{"post_id": "%s", "reaction": "like"}`, setup.postId),
			expectCode: 201,
		},
		{
			testName:   "incorrect post id",
			inputBody:  `{"post_id": "0", "reaction": "boom"}`,
			expectCode: 400,
		},
	}
	for _, tc := range testCases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts/reaction/create", bytes.NewBufferString(tc.inputBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+setup.token)
		s.router.ServeHTTP(w, req)
		s.Assert().Equal(tc.expectCode, w.Code)
		if tc.expectCode == 201 {
			var response struct {
				ReactionId string `json:"reaction_id"`
			}
			_ = json.Unmarshal(w.Body.Bytes(), &response)
			_, err := s.services.Reaction.GetReactionById(context.Background(), response.ReactionId)
			s.Assert().Equal(nil, err)
		}
	}
}
