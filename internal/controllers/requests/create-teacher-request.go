package requests

type CreateTeacherRequest struct {
	Fio         string `json:"fio"`
	PhoneNumber string `json:"phone_number"`
	AccountId   int    `json:"account_id"`
}
