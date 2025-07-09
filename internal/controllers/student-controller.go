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

type StudentController struct {
	readGroupUsecase                ReadGroupUsecase
	readAllStudentsUsecase          ReadAllStudentsUsecase
	readAllStudentsByGroupIdUsecase ReadAllStudentsByGroupIdUsecase
	readStudentUsecase              ReadStudentUsecase
	updateStudentUsecase            UpdateStudentUsecase
	deleteStudentUsecase            DeleteStudentUsecase
}

func NewStudentController(readGroupUsecase ReadGroupUsecase, readAllStudentsUsecase ReadAllStudentsUsecase, readAllStudentsByGroupIdUsecase ReadAllStudentsByGroupIdUsecase, readStudentUsecase ReadStudentUsecase, updateStudentUsecase UpdateStudentUsecase, deleteStudentUsecase DeleteStudentUsecase) StudentController {
	return StudentController{readGroupUsecase: readGroupUsecase, readAllStudentsUsecase: readAllStudentsUsecase, readAllStudentsByGroupIdUsecase: readAllStudentsByGroupIdUsecase, readStudentUsecase: readStudentUsecase, updateStudentUsecase: updateStudentUsecase, deleteStudentUsecase: deleteStudentUsecase}
}

// ReadAllStudents
// @Summary      Get all students
// @Description  Returns list of all students (admin only)
// @Tags         students
// @Security     BasicAuth
// @Produce      json
// @Success      200 {array} entities.Student
// @Failure      401 {object} object "Unauthorized"
// @Failure      403 {object} object "Forbidden"
// @Failure      500 {object} object "Internal server error"
// @Router       /api/read-all-students [get]
func (controller *StudentController) ReadAllStudents(c *gin.Context) {
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

	data, err := controller.readAllStudentsUsecase.ReadAllStudents(c)
	if err != nil {
		fmt.Println("failed to read students")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, data)
}

// ReadAllStudentsByGroupId
// @Summary      Get students by group ID
// @Description  Students of group (student sees own group, teacher sees own group, admin sees all)
// @Tags         students
// @Security     BasicAuth
// @Produce      json
// @Param        id query int true "Group ID"
// @Success      200 {object} usecases.ReadAllStudentsByGroupIdResponseDto
// @Failure      400 {object} object "Invalid group ID"
// @Failure      401 {object} object "Unauthorized"
// @Failure      403 {object} object "Forbidden"
// @Failure      500 {object} object "Internal server error"
// @Router       /api/read-all-students-by-group-id [get]
func (controller *StudentController) ReadAllStudentsByGroupId(c *gin.Context) {
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

	groupIdStr := c.Query("id")
	if groupIdStr == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	switch u := user.(type) {
	case entities.Student:
		if u.GroupId != groupId {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

	case entities.Teacher:
		group, err := controller.readGroupUsecase.ReadGroup(c, usecases.ReadGroupRequestDto{Id: groupId})
		if err != nil || group.Group.TeacherId != u.Id {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

	case entities.Admin:

	default:
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	data, err := controller.readAllStudentsByGroupIdUsecase.ReadAllStudentsByGroupId(
		c, usecases.ReadAllStudentsByGroupIdRequestDto{GroupId: groupId},
	)
	if err != nil {
		fmt.Println("failed to read students:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, data)
}

// ReadStudent
// @Summary      Get student by ID
// @Description  Returns student by ID. Accessible by student (self), teacher of group, or admin
// @Tags         students
// @Security     BasicAuth
// @Produce      json
// @Param        id query int true "Student ID"
// @Success      200 {object} usecases.ReadStudentResponseDto
// @Failure      400 {object} object "Invalid student ID"
// @Failure      401 {object} object "Unauthorized"
// @Failure      403 {object} object "Forbidden"
// @Failure      500 {object} object "Internal server error"
// @Router       /api/read-student [get]
func (controller *StudentController) ReadStudent(c *gin.Context) {
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

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	data, err := controller.readStudentUsecase.ReadStudent(c, usecases.ReadStudentRequestDto{Id: int(id)})
	if err != nil {
		fmt.Println("failed to read student:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	switch u := user.(type) {
	case entities.Student:
		if u.Id != int(id) || u.GroupId != data.Student.GroupId {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

	case entities.Teacher:
		group, err := controller.readGroupUsecase.ReadGroup(c, usecases.ReadGroupRequestDto{Id: data.Student.GroupId})
		if err != nil || group.Group.TeacherId != u.Id {
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

// UpdateStudent
// @Summary      Update student
// @Description  Update a student record.
//
//	Access allowed to:
//	- The student themselves
//	- Admins
//	- Teachers
//
// @Tags         students
// @Security     BasicAuth
// @Accept       json
// @Produce      json
// @Param        student body requests.UpdateStudentRequest true "Updated student info"
// @Success      200 {object} entities.Student
// @Failure      400 {object} object "Invalid request body"
// @Failure      403 {object} object "Access forbidden"
// @Failure      401 {object} object "Unauthorized"
// @Failure      500 {object} object "Internal server error"
// @Router       /api/update-student [put]
func (controller *StudentController) UpdateStudent(c *gin.Context) {
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

	req := requests.UpdateStudentRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	switch u := user.(type) {
	case entities.Student:
		if u.Id != req.Id {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

	case entities.Teacher:

	case entities.Admin:

	default:
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	data, err := controller.updateStudentUsecase.UpdateStudent(c, usecases.UpdateStudentRequestDto{Id: req.Id, Fio: req.Fio, PhoneNumber: req.PhoneNumber, GroupId: req.GroupId})
	if err != nil {
		fmt.Println("failed to update student")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, data)
}

// DeleteStudent
// @Summary      Delete student
// @Description  Delete student by ID (admin only)
// @Tags         students
// @Security     BasicAuth
// @Produce      json
// @Param        id query int true "Student ID"
// @Success      200
// @Failure      400 {object} object "Invalid student ID"
// @Failure      401 {object} object "Unauthorized"
// @Failure      403 {object} object "Forbidden"
// @Failure      500 {object} object "Internal server error"
// @Router       /api/delete-student [delete]
func (controller *StudentController) DeleteStudent(c *gin.Context) {
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

	err = controller.deleteStudentUsecase.DeleteStudent(c, usecases.DeleteStudentRequestDto{Id: int(id)})
	if err != nil {
		fmt.Println("failed to delete student:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
