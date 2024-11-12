package model

type Customer struct {
	ID    uint   `json:"id"`
	Name  string `json:"name" validate:"required"`
	CPF   string `json:"cpf" validate:"required"`
	Email string `json:"email" validate:"required"`
}
