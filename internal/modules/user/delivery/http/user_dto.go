package http

import "golab/internal/modules/user/usecase"

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=120"`
	Email string `json:"email" validate:"required,email,max=255"`
}

type PatchUserRequest struct {
	Name  *string `json:"name" validate:"omitempty,min=2,max=120"`
	Email *string `json:"email" validate:"omitempty,email,max=255"`
}

type UserResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func toResponse(out *usecase.UserOutput) UserResponse {
	return UserResponse{
		ID:        out.ID,
		Name:      out.Name,
		Email:     out.Email,
		CreatedAt: out.CreatedAt,
		UpdatedAt: out.UpdatedAt,
	}
}
