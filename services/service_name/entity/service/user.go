package service

import "awesomeProject/services/service_name/entity/repository"

type User struct {
	ID       string
	Username string
	Email    string
	Password string
}

func (u *User) ToRepository() *repository.UserPg {
	return &repository.UserPg{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}
}

func UserFromRepository(userPg *repository.UserPg) *User {
	return &User{
		ID:       userPg.ID,
		Username: userPg.Username,
		Email:    userPg.Email,
		Password: userPg.Password,
	}
}
