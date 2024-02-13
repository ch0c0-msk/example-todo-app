package handler

import (
	"errors"
	"net/http"

	todo "github.com/ch0c0-msk/example-todo-app"
	"github.com/gin-gonic/gin"
)

func (h *ApiHandler) signUp(c *gin.Context) {
	var user todo.User

	if err := c.BindJSON(&user); err != nil {
		newErrorResponse(c, errors.New("empty required fields"), http.StatusBadRequest)
		return
	}

	id, err := h.service.Authorization.CreateUser(user)
	if err != nil {
		newErrorResponse(c, errors.New("auth service error"), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type UserSignInData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *ApiHandler) signIn(c *gin.Context) {
	var userData UserSignInData

	if err := c.BindJSON(&userData); err != nil {
		newErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	token, err := h.service.Authorization.GenerateToken(userData.Username, userData.Password)
	if err != nil {
		newErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
