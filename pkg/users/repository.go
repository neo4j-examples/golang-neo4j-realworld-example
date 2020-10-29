package users

type UserRepository interface {
	RegisterUser(user User) error
}
