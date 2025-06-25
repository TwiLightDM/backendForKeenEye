package requests

type CreateStudentRequest struct {
	Fio         string `json:"fio"`
	GroupName   string `json:"group_name"`
	PhoneNumber string `json:"phone_number"`
}
