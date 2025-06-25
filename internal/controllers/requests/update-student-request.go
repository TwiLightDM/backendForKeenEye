package requests

type UpdateStudentRequest struct {
	Id          int    `json:"id"`
	Fio         string `json:"fio"`
	GroupName   string `json:"group_name"`
	PhoneNumber string `json:"phone_number"`
}
