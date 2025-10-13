package dtos

type RegisterRequest struct {
	Email string 		`json:"email"`
	Password string	`json:"password"`
	Name string			`json:"name"`
}

type LoginRequest struct {
	Email string		`json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID string 			`json:"id"`
	Email string 		`json:"email"`
	Name string			`json:"Name"`	
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
}