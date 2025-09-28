package model

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=5,max=30"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5,max=30"`
}
