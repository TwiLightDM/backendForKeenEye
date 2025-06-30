package entities

type Account struct {
	Id       int
	Login    string
	Password string
	Salt     string
}
