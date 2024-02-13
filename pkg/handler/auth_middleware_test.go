package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ch0c0-msk/example-todo-app/pkg/service"
	mock_service "github.com/ch0c0-msk/example-todo-app/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, token string)
	testCases := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockAuthorization, token string) {
				r.EXPECT().ParseToken(token).Return(1, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "1",
		},
		{
			name:                 "Empty auth header",
			headerName:           "",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"empty authorization header"}`,
		},
		{
			name:                 "Invalid auth header: invalid value",
			headerName:           "Authorization",
			headerValue:          "Bearr token",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"invalid authorization header"}`,
		},
		{
			name:                 "Invalid auth header: empty token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"invalid authorization header"}`,
		},
		{
			name:        "Parse Error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockAuthorization, token string) {
				r.EXPECT().ParseToken(token).Return(0, errors.New("some parsing error"))
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"parsing auth token error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			tc.mockBehavior(auth, tc.token)

			service := &service.Service{Authorization: auth}
			handler := ApiHandler{service}

			r := gin.New()
			r.GET("/identity", handler.userIdentity, func(c *gin.Context) {
				id, _ := c.Get(userCtx)
				c.String(http.StatusOK, "%d", id)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/identity", nil)
			req.Header.Set(tc.headerName, tc.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tc.expectedResponseBody)
		})
	}
}
