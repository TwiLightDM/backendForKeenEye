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

type TeacherController struct {
	readAllTeachersUsecase ReadAllTeachersUsecase
	readTeacherUsecase     ReadTeacherUsecase
	updateTeacherUsecase   UpdateTeacherUsecase
	deleteTeacherUsecase   DeleteTeacherUsecase
}

func NewTeacherController(readAllTeachersUsecase ReadAllTeachersUsecase, readTeacherUsecase ReadTeacherUsecase, updateTeacherUsecase UpdateTeacherUsecase, deleteTeacherUsecase DeleteTeacherUsecase) TeacherController {
	return TeacherController{readAllTeachersUsecase: readAllTeachersUsecase, readTeacherUsecase: readTeacherUsecase, updateTeacherUsecase: updateTeacherUsecase, deleteTeacherUsecase: deleteTeacherUsecase}
}

// ReadAllTeachers
// @Summary      Get all teachers
// @Description  Returns list of all teachers (admin only)
// @Tags         teachers
// @Security     BasicAuth
// @Produce      json
// @Success      200 {array} entities.Teacher
// @Failure      401 {object} object "Unauthorized"
// @Failure      403 {object} object "Forbidden"
// @Failure      500 {object} object "Internal server error"
// @Router       /api/read-all-teachers [get]
func (controller *TeacherController) ReadAllTeachers(c *gin.Context) {
	data, err := controller.readAllTeachersUsecase.ReadAllTeachers(c)
	if err != nil {
		fmt.Println("failed to read teachers")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, data)
}

// ReadTeacher
// @Summary      Get teacher by ID
// @Description  Get teacher by ID (teacher sees self, admin sees all, students forbidden)
// @Tags         teachers
// @Security     BasicAuth
// @Produce      json
// @Param        id query int true "Teacher ID"
// @Success      200 {object} usecases.ReadTeacherResponseDto
// @Failure      400 {object} object "Invalid teacher ID"
// @Failure      401 {object} object "Unauthorized"
// @Failure      403 {object} object "Forbidden"
// @Failure      500 {object} object "Internal server error"
// @Router       /api/read-teacher [get]
func (controller *TeacherController) ReadTeacher(c *gin.Context) {
	var user any
	var ok bool

	for _, key := range []string{"teacher", "admin"} {
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

	switch u := user.(type) {

	case entities.Teacher:
		if u.Id != int(id) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

	case entities.Admin:

	default:
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	data, err := controller.readTeacherUsecase.ReadTeacher(c, usecases.ReadTeacherRequestDto{Id: int(id)})
	if err != nil {
		fmt.Println("failed to read teacher:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, data)
}

// UpdateTeacher
// @Summary      Update teacher
// @Description  Update teacher info (teacher updates self, admin updates any)
// @Tags         teachers
// @Security     BasicAuth
// @Accept       json
// @Produce      json
// @Param        teacher body requests.UpdateTeacherRequest true "Updated teacher info"
// @Success      200 {object} entities.Teacher
// @Failure      400 {object} object "Invalid request"
// @Failure      401 {object} object "Unauthorized"
// @Failure      403 {object} object "Forbidden"
// @Failure      500 {object} object "Internal server error"
// @Router       /api/update-teacher [put]
func (controller *TeacherController) UpdateTeacher(c *gin.Context) {
	var (
		teacher entities.Teacher
		admin   entities.Admin
		ok      bool
	)

	if t, exists := c.Get("teacher"); exists {
		teacher, ok = t.(entities.Teacher)
	} else if a, exists := c.Get("admin"); exists {
		admin, ok = a.(entities.Admin)
	}

	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	req := requests.UpdateTeacherRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	switch {

	case teacher.Id > 0:
		if teacher.Id != req.Id {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

	case admin.Id > 0:

	default:
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	data, err := controller.updateTeacherUsecase.UpdateTeacher(c, usecases.UpdateTeacherRequestDto{Id: req.Id, Fio: req.Fio, PhoneNumber: req.PhoneNumber})
	if err != nil {
		fmt.Println("failed to update teacher")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, data)
}

// DeleteTeacher
// @Summary      Delete teacher
// @Description  Delete teacher by ID (admin only)
// @Tags         teachers
// @Security     BasicAuth
// @Produce      json
// @Param        id query int true "Teacher ID"
// @Success      200
// @Failure      400 {object} object "Invalid teacher ID"
// @Failure      401 {object} object "Unauthorized"
// @Failure      403 {object} object "Forbidden"
// @Failure      500 {object} object "Internal server error"
// @Router       /api/delete-teacher [delete]
func (controller *TeacherController) DeleteTeacher(c *gin.Context) {
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

	err = controller.deleteTeacherUsecase.DeleteTeacher(c, usecases.DeleteTeacherRequestDto{Id: int(id)})
	if err != nil {
		fmt.Println("failed to delete teacher:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
