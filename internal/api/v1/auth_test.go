package v1

import (
	"API_for_SN_go/internal/mocks/servicemocks"
	"API_for_SN_go/internal/model/pgmodel"
	"API_for_SN_go/internal/service"
	"API_for_SN_go/pkg/validator"
	"bytes"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthRouter_signUp(t *testing.T) {
	type args struct {
		ctx   context.Context
		input service.UserCreateInput
	}
	type MockBehaviour func(m *servicemocks.MockAuth, args args)

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
				input: service.UserCreateInput{
					Username:  "vasek",
					FirstName: "Vasya",
					LastName:  "Pupkin",
					Email:     "vasiliy@gmail.com",
					Password:  "1234",
				},
			},
			inputBody: `{"username": "vasek", "first_name": "Vasya", "last_name": "Pupkin", "email": "vasiliy@gmail.com", "password": "1234"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {
				m.EXPECT().CreateUser(args.ctx, args.input).Return(nil)
			},
			expectCode: 201,
			expectBody: "",
		},
		{
			testName:      "incorrect username",
			args:          args{ctx: context.Background()},
			inputBody:     `{"username": "Vsay@QD=-3", "first_name": "Vasya", "last_name": "Pupkin", "email": "vasiliy@gmail.com", "password": "1234"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {},
			expectCode:    400,
			expectBody:    `{"message":"field username can only consist of lower Latin characters, numbers and underscore symbol. Min length is 3, max: 32"}` + "\n",
		},
		{
			testName:      "incorrect email",
			args:          args{ctx: context.Background()},
			inputBody:     `{"username": "vasek", "first_name": "Vasya", "last_name": "Pupkin", "email": "asdasda3124p9-09as-d@gmail.com", "password": "1234"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {},
			expectCode:    400,
			expectBody:    `{"message":"field email is incorrect. Make sure that you entered the email correctly and it exists"}` + "\n",
		},
		{
			testName:      "too long username",
			args:          args{ctx: context.Background()},
			inputBody:     `{"username": "looooooooooooooooooooooooongvasya", "first_name": "Vasya", "last_name": "Pupkin", "email": "vasiliy@gmail.com", "password": "1234"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {},
			expectCode:    400,
			expectBody:    `{"message":"field username can only consist of lower Latin characters, numbers and underscore symbol. Min length is 3, max: 32"}` + "\n",
		},
		{
			testName:      "too short username",
			args:          args{ctx: context.Background()},
			inputBody:     `{"username": "v", "first_name": "Vasya", "last_name": "Pupkin", "email": "vasiliy@gmail.com", "password": "1234"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {},
			expectCode:    400,
			expectBody:    `{"message":"field username can only consist of lower Latin characters, numbers and underscore symbol. Min length is 3, max: 32"}` + "\n",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			auth := servicemocks.NewMockAuth(ctrl)
			tc.mockBehaviour(auth, tc.args)
			services := &service.Services{Auth: auth}

			e := echo.New()
			e.Validator, _ = validator.NewValidator()
			newAuthRouter(e.Group("/auth"), services.Auth)

			// create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBufferString(tc.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			// execute request
			e.ServeHTTP(w, req)

			// check response
			assert.Equal(t, tc.expectCode, w.Code)
			assert.Equal(t, tc.expectBody, w.Body.String())
		})
	}
}

func (s *APITestSuite) Test_authRouter_signUp() {
	testCases := []struct {
		testName   string
		input      pgmodel.User
		inputBody  string
		expectCode int
		expectBody string
	}{
		{
			testName: "Correct test",
			input: pgmodel.User{
				Username:  "vasek",
				FirstName: "Vasya",
				LastName:  "Pupkin",
				Email:     "vasiliy@gmail.com",
				Password:  "1234",
			},
			inputBody:  `{"username": "vasek", "first_name": "Vasya", "last_name": "Pupkin", "email": "vasiliy@gmail.com", "password": "1234"}`,
			expectCode: 201,
			expectBody: "",
		},
		{
			testName: "Repeated username test",
			input: pgmodel.User{
				Id:        0,
				Username:  "vasek",
				FirstName: "Petya",
				LastName:  "Petrov",
				Email:     "petrov@gmail.com",
				Password:  "1234",
			},
			inputBody:  `{"username": "vasek", "first_name": "Petya", "last_name": "Petrov", "email": "petrov@gmail.com", "password": "1234"}`,
			expectCode: 400,
			expectBody: `{"message":"user already exists"}` + "\n",
		},
	}
	for _, tc := range testCases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBufferString(tc.inputBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		s.router.ServeHTTP(w, req)

		s.Assert().Equal(tc.expectCode, w.Code)
		s.Assert().Equal(tc.expectBody, w.Body.String())

		if tc.expectCode == 201 {
			u, err := s.services.User.GetUserByUsername(context.Background(), tc.input.Username)
			if err != nil {
				panic(err)
			}
			tc.input.Id = u.Id
			tc.input.Password = u.Password
			s.Assert().Equal(tc.input, u)
		}
	}
}

func (s *APITestSuite) Test_authRouter_signIn() {
	if err := s.services.Auth.CreateUser(context.Background(), service.UserCreateInput{
		Username:  "vasek",
		FirstName: "Vasya",
		LastName:  "Pupkin",
		Email:     "test",
		Password:  "1234",
	}); err != nil {
		panic(err)
	}
	testCases := []struct {
		testName   string
		inputBody  string
		expectCode int
		expectBody string
	}{
		{
			testName:   "correct test",
			inputBody:  `{"username": "vasek", "password": "1234"}`,
			expectCode: 200,
		},
		{
			testName:   "incorrect user password",
			inputBody:  `{"username": "vasek", "password": "my pass 1234567890"}`,
			expectCode: 403,
			expectBody: `{"message":"incorrect user password"}` + "\n",
		},
	}
	for _, tc := range testCases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/sign-in", bytes.NewBufferString(tc.inputBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		s.router.ServeHTTP(w, req)

		s.Assert().Equal(tc.expectCode, w.Code)

		if tc.expectCode == 200 {
			token, err := s.redis.Pool.Get(context.Background(), "jwt:vasek").Result()
			s.Assert().Equal(nil, err)
			s.Assert().NotEqual("", token)
		} else {
			s.Assert().Equal(tc.expectBody, w.Body.String())
		}
	}
}
