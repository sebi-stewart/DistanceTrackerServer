package auth

import "fmt"

func ValidateRegistration(user UserRegister) error {
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}
	if user.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if user.Password == "" {
		return fmt.Errorf("password is required")
	}
	if user.ConfirmPassword == "" {
		return fmt.Errorf("confirm password is required")
	}
	if user.Password != user.ConfirmPassword {
		return fmt.Errorf("passwords do not match")
	}
	return nil
}
