package requests

type UpdateStudentRequest struct {
	Id          int    `json:"id"`
	Fio         string `json:"fio"`
	PhoneNumber string `json:"phone_number"`
	GroupId     int    `json:"group_id"`
}
