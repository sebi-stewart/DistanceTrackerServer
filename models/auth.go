package models

import (
	"fmt"
	"github.com/google/uuid"
)

type UserRegister struct {
	Email           string `json:"email"`
	FirstName       string `json:"first_name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (u *UserRegister) ToString() string {
	return fmt.Sprintf("{email: %s,\tfirst_name: %s,\tpassword: %s,\tconfirm_password: %s}",
		u.Email, u.FirstName, u.Password, u.ConfirmPassword,
	)
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AccountLink struct {
	PairUUID uuid.UUID `json:"pair_uuid"`
}

func (a *AccountLink) ToString() string {
	return fmt.Sprintf("{pair_uuid: %s}",
		a.PairUUID.String(),
	)
}
