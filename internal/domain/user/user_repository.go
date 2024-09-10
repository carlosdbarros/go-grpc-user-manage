package user

type UserRepository interface {
	AddUser(user *User) (*User, error)
	FindUserByEmail(email string) (*User, error)
	FindAllUsers() ([]*User, error)
}
