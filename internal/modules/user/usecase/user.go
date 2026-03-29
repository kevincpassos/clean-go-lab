package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"golab/internal/modules/user/domain"
	"golab/internal/modules/user/usecase/ports"
)

type UserUseCase struct {
	Repo     ports.UserRepository
	Mailer   ports.Mailer
	EventBus ports.EventBus
}

func NewUserUseCase(
	repo ports.UserRepository,
	mailer ports.Mailer,
	eventBus ports.EventBus,
) *UserUseCase {
	return &UserUseCase{
		Repo:     repo,
		Mailer:   mailer,
		EventBus: eventBus,
	}
}

func (uc *UserUseCase) Create(ctx context.Context, input CreateUserInput) (*UserOutput, error) {
	user, err := domain.NewUser(input.Name, input.Email)
	if err != nil {
		return nil, err
	}

	created, err := uc.Repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(AccountActivationEmailMessage{
		UserID: created.ID,
		Email:  created.Email,
		Name:   created.Name,
	})
	if err != nil {
		return nil, ErrCreateActivationPayloadMarshal
	}

	if err := uc.EventBus.PublishAccountActivationEmail(ctx, payload); err != nil {
		return nil, ErrCreateActivationEventPublish
	}

	return ToUserOutput(created), nil
}

func (uc *UserUseCase) Patch(ctx context.Context, input PatchUserInput) (*UserOutput, error) {
	if err := validatePatchInput(input); err != nil {
		return nil, err
	}

	updated, err := uc.Repo.Update(ctx, ports.PatchUserParams{
		ID:    input.ID,
		Name:  input.Name,
		Email: input.Email,
	})
	if err != nil {
		return nil, err
	}

	return ToUserOutput(updated), nil
}

func (uc *UserUseCase) Delete(ctx context.Context, input DeleteUserInput) error {
	if input.ID <= 0 {
		return ErrInvalidUserID
	}

	if err := uc.Repo.Delete(ctx, input.ID); err != nil {
		return err
	}

	return nil
}

func (uc *UserUseCase) SendActivationEmail(ctx context.Context, input SendActivationEmailInput) error {
	if err := uc.Mailer.SendActivationEmail(ctx, input.Email, input.Name); err != nil {
		return errors.Join(ErrSendActivationEmail, err)
	}
	return nil
}
