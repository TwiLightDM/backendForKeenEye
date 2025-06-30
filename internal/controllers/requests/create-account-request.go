package requests

type CreateAccountRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
