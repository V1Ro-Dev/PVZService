package models

type User struct {
	Email    string
	Password string
	Salt     string
	Role     string
	Id       string
}

type LoginData struct {
	Email    string
	Password string
}
