package models

import (
	"errors"
	"strings"
)

const (
	DefaultRoleId = 9
	Activated     = 1
	deactivated   = 0

	BcryptCose = 10
)

var ErrAccountDuplicated = errors.New("This account is exist.")

var ErrWrongPassword = errors.New("Wrong password.")

var ErrNoAccount = errors.New("This account is not exist.")

type Account struct {
	Account  string `gorm:"primary_key" json:"account"`
	Password string `json:"-"`
	RoleId   int    `json:"roleId"`
	Role     Role   `gorm:"ForeignKey:RoleId;save_associations:false"`
	Status   int

	BaseModel
}

type AccountForm struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

func (af AccountForm) CheckField() error {
	errMsg := make([]string, 0)

	if len(af.Account) == 0 {
		errMsg = append(errMsg, "Account should not be empty")
	}

	if len(af.Password) == 0 {
		errMsg = append(errMsg, "Password should not be empty")
	}

	if len(errMsg) > 0 {
		return errors.New(strings.Join(errMsg, ", ") + ".")
	} else {
		return nil
	}
}
