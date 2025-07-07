package controllers

import (
	"backendForKeenEye/internal/controllers/requests"
	"backendForKeenEye/internal/entities"
	"backendForKeenEye/internal/usecases"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AdminController struct {
	createAdminUsecase CreateAdminUsecase
	readAdminUsecase   ReadAdminUsecase
	updateAdminUsecase UpdateAdminUsecase
	deleteAdminUsecase DeleteAdminUsecase
}

func NewAdminController(createAdminUsecase CreateAdminUsecase, readAdminUsecase ReadAdminUsecase, updateAdminUsecase UpdateAdminUsecase, deleteAdminUsecase DeleteAdminUsecase) AdminController {
	return AdminController{createAdminUsecase: createAdminUsecase, readAdminUsecase: readAdminUsecase, updateAdminUsecase: updateAdminUsecase, deleteAdminUsecase: deleteAdminUsecase}
}

// CreateAdmin
// @Summary      Create admin
// @Description  Create a new admin (admin only)
// @Tags         admins
// @Security     BasicAuth
// @Accept       json
// @Produce      json
// @Param        admin body requests.CreateAdminRequest true "Admin info"
// @Success      201 {object} entities.Admin
// @Failure      400 {object} object "Invalid request"
// @Failure      401 {object} object "Unauthorized"
// @Failure      403 {object} object "Forbidden"
// @Failure      500 {object} object "Internal server error"
// @Router       /api/create-admin [post]
func (controller *AdminController) CreateAdmin(c *gin.Context) {
	accountRaw, exists := c.Get("account")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	account, ok := accountRaw.(entities.Account)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if account.Role != "admin" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	req := requests.CreateAdminRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	data, err := controller.createAdminUsecase.CreateAdmin(c, usecases.CreateAdminRequestDto{Fio: req.Fio, PhoneNumber: req.PhoneNumber, AccountId: req.AccountId})
	if err != nil {
		fmt.Println("failed to create admin", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, data)
}

// ReadAdmin
// @Summary      Get admin by ID
// @Description  Get admin by ID (admin only)
// @Tags         admins
// @Security     BasicAuth
// @Produce      json
// @Param        id query int true "Admin ID"
// @Success      200 {object} usecases.ReadAdminResponseDto
// @Failure      400 {object} object "Invalid admin ID"
// @Failure      401 {object} object "Unauthorized"
// @Failure      403 {object} object "Forbidden"
// @Failure      500 {object} object "Internal server error"
// @Router       /api/read-admin [get]
func (controller *AdminController) ReadAdmin(c *gin.Context) {
	accountRaw, exists := c.Get("account")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	account, ok := accountRaw.(entities.Account)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if account.Role != "admin" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	idStr := c.Query("id")
	if idStr == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	data, err := controller.readAdminUsecase.ReadAdmin(c, usecases.ReadAdminRequestDto{Id: int(id)})
	if err != nil {
		fmt.Println("failed to read admin:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, data)
}

// UpdateAdmin
// @Summary      Update admin
// @Description  Update admin info (admin only)
// @Tags         admins
// @Security     BasicAuth
// @Accept       json
// @Produce      json
// @Param        admin body requests.UpdateAdminRequest true "Updated admin info"
// @Success      200 {object} entities.Admin
// @Failure      400 {object} object "Invalid request"
// @Failure      401 {object} object "Unauthorized"
// @Failure      403 {object} object "Forbidden"
// @Failure      500 {object} object "Internal server error"
// @Router       /api/update-admin [put]
func (controller *AdminController) UpdateAdmin(c *gin.Context) {
	accountRaw, exists := c.Get("account")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	account, ok := accountRaw.(entities.Account)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if account.Role != "admin" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	req := requests.UpdateAdminRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	data, err := controller.updateAdminUsecase.UpdateAdmin(c, usecases.UpdateAdminRequestDto{Id: req.Id, Fio: req.Fio, PhoneNumber: req.PhoneNumber, AccountId: req.AccountId})
	if err != nil {
		fmt.Println("failed to update admin")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, data)
}

// DeleteAdmin
// @Summary      Delete admin
// @Description  Delete admin by ID (admin only)
// @Tags         admins
// @Security     BasicAuth
// @Produce      json
// @Param        id query int true "Admin ID"
// @Success      200
// @Failure      400 {object} object "Invalid admin ID"
// @Failure      401 {object} object "Unauthorized"
// @Failure      403 {object} object "Forbidden"
// @Failure      500 {object} object "Internal server error"
// @Router       /api/delete-admin [delete]
func (controller *AdminController) DeleteAdmin(c *gin.Context) {
	accountRaw, exists := c.Get("account")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	account, ok := accountRaw.(entities.Account)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if account.Role != "admin" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	idStr := c.Query("id")
	if idStr == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = controller.deleteAdminUsecase.DeleteAdmin(c, usecases.DeleteAdminRequestDto{Id: int(id)})
	if err != nil {
		fmt.Println("failed to delete admin:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
