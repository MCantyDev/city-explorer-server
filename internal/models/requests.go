package models

type SignupRequest struct {
	FirstName string `json:"first_name" binding:"required,min=1,max=32"`
	LastName  string `json:"last_name" binding:"required,min=1,max=32"`
	Username  string `json:"username" binding:"required,min=3,max=32"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6"`
}
