package controllers

import (
	"backendForKeenEye/internal/controllers/requests"
	"backendForKeenEye/internal/usecases"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type StudentController struct {
	createStudentUsecase   CreateStudentUsecase
	readAllStudentsUsecase ReadAllStudentsUsecase
	readStudentUsecase     ReadStudentUsecase
	updateStudentUsecase   UpdateStudentUsecase
	deleteStudentUsecase   DeleteStudentUsecase
}

func NewStudentController(createStudentUsecase CreateStudentUsecase, readAllStudentsUsecase ReadAllStudentsUsecase, readStudentUsecase ReadStudentUsecase, updateStudentUsecase UpdateStudentUsecase, deleteStudentUsecase DeleteStudentUsecase) StudentController {
	return StudentController{createStudentUsecase: createStudentUsecase, readAllStudentsUsecase: readAllStudentsUsecase, readStudentUsecase: readStudentUsecase, updateStudentUsecase: updateStudentUsecase, deleteStudentUsecase: deleteStudentUsecase}
}

// CreateStudent
// @Summary      Create student
// @Description  Create a new student
// @Tags         students
// @Accept       json
// @Produce      json
// @Param        student body requests.CreateStudentRequest true "Student info"
// @Success      201 {object} entities.Student
// @Failure      400 {object} object
// @Failure      500 {object} object
// @Router       /api/create-student [post]
func (controller *StudentController) CreateStudent(c *gin.Context) {
	req := requests.CreateStudentRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	data, err := controller.createStudentUsecase.CreateStudent(c, usecases.CreateStudentRequestDto{Fio: req.Fio, GroupName: req.GroupName, PhoneNumber: req.PhoneNumber})
	if err != nil {
		fmt.Println("failed to create student", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, data)
}

// ReadAllStudents
// @Summary      Get all students
// @Description  Returns list of all students
// @Tags         students
// @Produce      json
// @Success      200 {array} entities.Student
// @Failure      500 {object} object
// @Router       /api/read-all-students [get]
func (controller *StudentController) ReadAllStudents(c *gin.Context) {
	data, err := controller.readAllStudentsUsecase.ReadAllStudents(c)
	if err != nil {
		fmt.Println("failed to read students")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, data)
}

// ReadStudent
// @Summary      Get student by ID
// @Description  Returns a student by ID
// @Tags         students
// @Produce      json
// @Param        id query int true "Student ID"
// @Success      200 {object} entities.Student
// @Failure      400 {object} object
// @Failure      500 {object} object
// @Router       /api/read-student [get]
func (controller *StudentController) ReadStudent(c *gin.Context) {
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

	c.JSON(http.StatusOK, data)
}

// UpdateStudent
// @Summary      Update student
// @Description  Update student info
// @Tags         students
// @Accept       json
// @Produce      json
// @Param        student body requests.UpdateStudentRequest true "Updated student info"
// @Success      200 {object} entities.Student
// @Failure      400 {object} object
// @Failure      500 {object} object
// @Router       /api/update-student [put]
func (controller *StudentController) UpdateStudent(c *gin.Context) {
	req := requests.UpdateStudentRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	data, err := controller.updateStudentUsecase.UpdateStudent(c, usecases.UpdateStudentRequestDto{Id: req.Id, Fio: req.Fio, GroupName: req.GroupName, PhoneNumber: req.PhoneNumber})
	if err != nil {
		fmt.Println("failed to update student")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, data)
}

// DeleteStudent
// @Summary      Delete student
// @Description  Delete student by ID
// @Tags         students
// @Produce      json
// @Param        id query int true "Student ID"
// @Success      200
// @Failure      400 {object} object
// @Failure      500 {object} object
// @Router       /api/delete-student [delete]
func (controller *StudentController) DeleteStudent(c *gin.Context) {
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
