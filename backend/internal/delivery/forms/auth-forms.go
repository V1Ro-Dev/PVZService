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
	Id    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func ToSignUpOut(user models.User) SignInFormOut {
	return SignInFormOut{
		Id:    user.Id,
		Email: user.Email,
		Role:  user.Role,
	}
}

type LogInFormIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
