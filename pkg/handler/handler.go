package handler

import (
	"github.com/ch0c0-msk/example-todo-app/pkg/service"
	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
	service *service.Service
}

func (h *ApiHandler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	auth.POST("/sign-up", h.signUp)
	auth.POST("/sign-in", h.signIn)

	return router
}

func NewHandler(service *service.Service) *ApiHandler {
	return &ApiHandler{service: service}
}
