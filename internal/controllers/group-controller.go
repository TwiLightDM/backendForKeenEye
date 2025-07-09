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

type GroupController struct {
	createGroupUsecase   CreateGroupUsecase
	readAllGroupsUsecase ReadAllGroupsUsecase
	readGroupUsecase     ReadGroupUsecase
	updateGroupUsecase   UpdateGroupUsecase
	deleteGroupUsecase   DeleteGroupUsecase
}

func NewGroupController(createGroupUsecase CreateGroupUsecase, readAllGroupsUsecase ReadAllGroupsUsecase, readGroupUsecase ReadGroupUsecase, updateGroupUsecase UpdateGroupUsecase, deleteGroupUsecase DeleteGroupUsecase) GroupController {
	return GroupController{createGroupUsecase: createGroupUsecase, readAllGroupsUsecase: readAllGroupsUsecase, readGroupUsecase: readGroupUsecase, updateGroupUsecase: updateGroupUsecase, deleteGroupUsecase: deleteGroupUsecase}
}

// CreateGroup
// @Summary      Create group
// @Description  Create a new group (admin only)
// @Tags         groups
// @Security     BasicAuth
// @Accept       json
// @Produce      json
// @Param        group body requests.CreateGroupRequest true "Group info"
// @Success      201 {object} entities.Group
// @Failure      400 {object} object "Invalid request"
// @Failure      401 {object} object "Unauthorized"
// @Failure      403 {object} object "Forbidden"
// @Failure      500 {object} object "Internal server error"
// @Router       /api/create-group [post]
func (controller *GroupController) CreateGroup(c *gin.Context) {
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

	req := requests.CreateGroupRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	data, err := controller.createGroupUsecase.CreateGroup(c, usecases.CreateGroupRequestDto{Name: req.Name, TeacherId: req.TeacherId})
	if err != nil {
		fmt.Println("failed to create group", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, data)
}

// ReadAllGroups
// @Summary      Get all groups
// @Description  Get list of all groups (admin only)
// @Tags         groups
// @Security     BasicAuth
// @Produce      json
// @Success      200 {object} usecases.ReadAllGroupsResponseDto
// @Failure      401 {object} object "Unauthorized"
// @Failure      403 {object} object "Forbidden"
// @Failure      500 {object} object "Internal server error"
// @Router       /api/read-all-groups [get]
func (controller *GroupController) ReadAllGroups(c *gin.Context) {
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

	data, err := controller.readAllGroupsUsecase.ReadAllGroups(c)
	if err != nil {
		fmt.Println("failed to read groups")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, data)
}

// ReadGroup
// @Summary      Get group by ID
// @Description  Get group by ID (student, teacher or admin). Students and teachers are allowed to access only their group.
// @Tags         groups
// @Security     BasicAuth
// @Produce      json
// @Param        id query int true "Group ID"
// @Success      200 {object} usecases.ReadGroupResponseDto
// @Failure      400 {object} object "Invalid group ID"
// @Failure      401 {object} object "Unauthorized"
// @Failure      403 {object} object "Forbidden"
// @Failure      500 {object} object "Internal server error"
// @Router       /api/read-group [get]
func (controller *GroupController) ReadGroup(c *gin.Context) {
	var user any
	var ok bool

	for _, key := range []string{"student", "teacher", "admin"} {
		if userRaw, exists := c.Get(key); exists {
			user = userRaw
			ok = true
			break
		}
	}

	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	idStr := c.Query("id")
	if idStr == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	data, err := controller.readGroupUsecase.ReadGroup(c, usecases.ReadGroupRequestDto{Id: id})
	if err != nil {
		fmt.Println("failed to read group:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	switch u := user.(type) {
	case entities.Student:
		if u.GroupId != id {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

	case entities.Teacher:
		if data.Group.TeacherId != u.Id {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

	case entities.Admin:

	default:
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, data)
}

// UpdateGroup
// @Summary      Update group
// @Description  Update group info (admin only)
// @Tags         groups
// @Security     BasicAuth
// @Accept       json
// @Produce      json
// @Param        group body requests.UpdateGroupRequest true "Updated group info"
// @Success      200 {object} entities.Group
// @Failure      400 {object} object "Invalid request"
// @Failure      401 {object} object "Unauthorized"
// @Failure      403 {object} object "Forbidden"
// @Failure      500 {object} object "Internal server error"
// @Router       /api/update-group [put]
func (controller *GroupController) UpdateGroup(c *gin.Context) {
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

	req := requests.UpdateGroupRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	data, err := controller.updateGroupUsecase.UpdateGroup(c, usecases.UpdateGroupRequestDto{Id: req.Id, Name: req.Name, TeacherId: req.TeacherId})
	if err != nil {
		fmt.Println("failed to update group")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, data)
}

// DeleteGroup
// @Summary      Delete group
// @Description  Delete group by ID (admin only)
// @Tags         groups
// @Security     BasicAuth
// @Produce      json
// @Param        id query int true "Group ID"
// @Success      200
// @Failure      400 {object} object "Invalid group ID"
// @Failure      401 {object} object "Unauthorized"
// @Failure      403 {object} object "Forbidden"
// @Failure      500 {object} object "Internal server error"
// @Router       /api/delete-group [delete]
func (controller *GroupController) DeleteGroup(c *gin.Context) {
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

	err = controller.deleteGroupUsecase.DeleteGroup(c, usecases.DeleteGroupRequestDto{Id: int(id)})
	if err != nil {
		fmt.Println("failed to delete group:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
