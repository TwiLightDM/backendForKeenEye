package controllers

import (
	"backendForKeenEye/internal/controllers/requests"
	"backendForKeenEye/internal/usecases"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AccountController struct {
	createAccountUsecase CreateAccountUsecase
}

func NewAccountController(createAccountUsecase CreateAccountUsecase) AccountController {
	return AccountController{createAccountUsecase: createAccountUsecase}
}

// CreateAccount
// @Summary      Create account
// @Description  Register a new user account
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        account body requests.CreateAccountRequest true "Account credentials"
// @Success      201 {object} entities.Account
// @Failure      400 {object} object
// @Failure      500 {object} object
// @Router       /api/create-account [post]
func (controller *AccountController) CreateAccount(c *gin.Context) {
	req := requests.CreateAccountRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	data, err := controller.createAccountUsecase.CreateAccount(c, usecases.CreateAccountRequestDto{Login: req.Login, Password: req.Password, Role: req.Role})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, data)
}
