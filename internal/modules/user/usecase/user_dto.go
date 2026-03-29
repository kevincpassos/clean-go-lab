package usecase

import "golab/internal/modules/user/domain"

type CreateUserInput struct {
	Name  string
	Email string
}

type PatchUserInput struct {
	ID    int64
	Name  *string
	Email *string
}

type DeleteUserInput struct {
	ID int64
}

type SendActivationEmailInput struct {
	Email string
	Name  string
}

// AccountActivationEmailMessage is the payload published to (and consumed from) the activation e-mail queue.
type AccountActivationEmailMessage struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type UserOutput struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt string
	UpdatedAt string
}

func ToUserOutput(user *domain.User) *UserOutput {
	return &UserOutput{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
