package usecase

import "golab/internal/modules/user/domain"

func validatePatchInput(input PatchUserInput) error {
	if input.ID <= 0 {
		return ErrInvalidUserID
	}

	if input.Name == nil && input.Email == nil {
		return ErrNoFieldsToUpdate
	}

	if input.Name != nil {
		if err := domain.ValidateName(input.Name); err != nil {
			return err
		}
	}

	if input.Email != nil {
		if err := domain.ValidateEmail(input.Email); err != nil {
			return err
		}
	}

	return nil
}
