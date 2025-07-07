package requests

type UpdateTeacherRequest struct {
	Id          int    `json:"id"`
	Fio         string `json:"fio"`
	PhoneNumber string `json:"phone_number"`
	AccountId   int    `json:"account_id"`
}
