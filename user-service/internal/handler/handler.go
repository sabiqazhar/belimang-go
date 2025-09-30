package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	errors2 "github.com/sabiqazhar/belimang-go/user-service/errors"
	"github.com/sabiqazhar/belimang-go/user-service/internal/model"
	"github.com/sabiqazhar/belimang-go/user-service/internal/service"
)

type UserHandler struct {
	userService service.UserService
	engine      *gin.Engine
}

func NewHandler(engine *gin.Engine, userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		engine:      engine,
	}
}

func (h *UserHandler) RegisterRoutes() {
	adminRoutes := h.engine.Group("/admin")
	adminRoutes.POST("/register", h.AdminRegister)
	adminRoutes.POST("/login", h.UserLogin)

	userRoutes := h.engine.Group("/users")
	userRoutes.POST("/register", h.UserRegister)
	userRoutes.POST("/login", h.UserLogin)
}

func (h *UserHandler) AdminRegister(c *gin.Context) {
	var req model.UserRegisterRequest
	ctx := c.Request.Context()
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := h.userService.CreateUser(
		ctx,
		req,
		true,
	)

	if err != nil {
		if errors.Is(err, errors2.ErrEmailAlreadyExists) || errors.Is(err, errors2.ErrAdminEmailExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create admin user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func (h *UserHandler) UserLogin(c *gin.Context) {
	var req model.UserLoginRequest
	ctx := c.Request.Context()
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	token, err := h.userService.UserLogin(ctx, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) UserRegister(c *gin.Context) {
	var req model.UserRegisterRequest
	ctx := c.Request.Context()
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := h.userService.CreateUser(
		ctx,
		req,
		false,
	)

	if err != nil {
		if errors.Is(err, errors2.ErrEmailAlreadyExists) || errors.Is(err, errors2.ErrAdminEmailExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"token": token})
}
