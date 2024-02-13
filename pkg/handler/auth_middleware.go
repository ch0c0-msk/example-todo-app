package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *ApiHandler) userIdentity(c *gin.Context) {

	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, errors.New("empty authorization header"), http.StatusUnauthorized)
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" || len(headerParts[1]) == 0 {
		newErrorResponse(c, errors.New("invalid authorization header"), http.StatusUnauthorized)
		return
	}

	userId, err := h.service.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, errors.New("parsing auth token error"), http.StatusUnauthorized)
		return
	}

	c.Set(userCtx, userId)
}
