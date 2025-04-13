package dto

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=5"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	Role     string `json:"role"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=5"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name"`
	Role     string `json:"role" binding:"required,oneof=admin warehouse_admin warehouse_manager"`
}

type UpdateUserRequest struct {
	Username *string `json:"username" binding:"omitempty,min=5"`
	FullName *string `json:"full_name" binding:"omitempty"`
	Password *string `json:"password" binding:"omitempty,min=8"`
	Role     *string `json:"role" binding:"omitempty,oneof=admin warehouse_admin warehouse_manager"`
}

type UserListResponse struct {
	Data       []UserResponse `json:"data"`
	Pagination Pagination     `json:"pagination"`
}
