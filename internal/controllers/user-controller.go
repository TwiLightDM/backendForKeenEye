package controllers

import (
	"backendForKeenEye/internal/controllers/requests"
	"backendForKeenEye/internal/entities"
	"backendForKeenEye/internal/usecases"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	createUserUsecase CreateUserUsecase
}

func NewUserController(createUserUsecase CreateUserUsecase) UserController {
	return UserController{createUserUsecase: createUserUsecase}
}

// CreateUser
// @Summary      Create user
// @Description  Create a new user (admin only)
// @Tags         users
// @Security     BasicAuth
// @Accept       json
// @Produce      json
// @Param        user body requests.CreateUserRequest true "User info"
// @Success      201 {object} entities.User
// @Failure      400 {object} object "Invalid request"
// @Failure      401 {object} object "Unauthorized"
// @Failure      403 {object} object "Forbidden"
// @Failure      500 {object} object "Internal server error"
// @Router       /api/create-user [post]
func (controller *UserController) CreateUser(c *gin.Context) {
	userRaw, exists := c.Get("user")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, ok := userRaw.(entities.User)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if user.Role != "admin" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	req := requests.CreateUserRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	data, err := controller.createUserUsecase.CreateUser(c, usecases.CreateUserRequestDto{Login: req.Login, Password: req.Password, Role: req.Role})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, data)
}
