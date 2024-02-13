package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	todo "github.com/ch0c0-msk/example-todo-app"
	"github.com/ch0c0-msk/example-todo-app/pkg/service"
	mock_service "github.com/ch0c0-msk/example-todo-app/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user todo.User)

	testCases := []struct {
		name                string
		inputBody           string
		inputUser           todo.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name" : "test-user-name", "username" : "test-user", "password" : "password"}`,
			inputUser: todo.User{
				Name:     "test-user-name",
				Username: "test-user",
				Password: "password",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:                "NOT OK: Empty required fields",
			inputBody:           `{"username" : "test-user", "password" : "password"}`,
			mockBehavior:        func(s *mock_service.MockAuthorization, user todo.User) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"error":"empty required fields"}`,
		},
		{
			name:      "NOT OK: Service error",
			inputBody: `{"name" : "test-user-name", "username" : "test-user", "password" : "password"}`,
			inputUser: todo.User{
				Name:     "test-user-name",
				Username: "test-user",
				Password: "password",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(0, errors.New("some service error"))
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: `{"error":"auth service error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			auth := mock_service.NewMockAuthorization(c)
			tc.mockBehavior(auth, tc.inputUser)
			service := &service.Service{Authorization: auth}
			handler := NewHandler(service)

			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/sign-up", bytes.NewBufferString(tc.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedRequestBody, w.Body.String())
		})
	}
}
