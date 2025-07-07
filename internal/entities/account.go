package entities

var allowedRoles = []string{"admin", "student", "teacher"}

type Account struct {
	Id       int
	Login    string
	Password string
	Salt     string
	Role     string
}

func (a Account) Validate() (bool, error) {
	if validateRole(a.Role) == false {
		return false, InvalidRoleError
	}
	return true, nil
}

func validateRole(role string) bool {
	for _, r := range allowedRoles {
		if role == r {
			return true
		}
	}
	return false
}
