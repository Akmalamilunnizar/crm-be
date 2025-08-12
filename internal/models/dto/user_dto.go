package dto

import "skripsi-be/internal/models/entities"

type UserDTO struct {
	ID      string        `json:"id"`
	Name    string        `json:"name" `
	Email   string        `json:"email" `
	Role    entities.Role `json:"role" `
	LogoUrl string        `json:"logo_url" `
	Phone   string        `json:"phone" `
}

func ModelToUser(user entities.User) *UserDTO {
	return &UserDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}
