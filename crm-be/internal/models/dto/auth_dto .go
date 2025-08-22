package dto

import "skripsi-be/internal/models/entities"

type LoginAuthDTO struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

func ModelToLoginAuth(token string, user entities.User) *LoginAuthDTO {
	return &LoginAuthDTO{
		Token: token,
		User: UserDTO{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}
}

type CustomerLoginAuthDTO struct {
	Token    string      `json:"token"`
	Customer CustomerDTO `json:"user"`
}

func ModelToCustomerLoginAuth(token string, user CustomerDTO) *CustomerLoginAuthDTO {
	return &CustomerLoginAuthDTO{
		Token:    token,
		Customer: user,
	}
}
