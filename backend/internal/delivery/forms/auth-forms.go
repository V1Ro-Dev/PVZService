package forms

import "pvz/internal/models"

type DummyLoginForm struct {
	Role string `json:"role"`
}

type SignUpFormIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type SignInFormOut struct {
	id    string
	email string
	role  string
}

func ToSignUpOut(user models.User) SignInFormOut {
	return SignInFormOut{
		id:    user.Id,
		email: user.Email,
		role:  user.Role,
	}
}
