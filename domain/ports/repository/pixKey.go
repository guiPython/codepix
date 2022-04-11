package repository

import (
	"fmt"

	"github.com/guiPython/codepix/domain/model"
)

type IPixKeyRepository interface {
	RegisterPixKey(pixKey *model.PixKey) (*model.PixKey, *ErrPixKeyRepository)

	FindPixKeyByKind(key string, kind model.Kind) (*model.PixKey, *ErrPixKeyRepository)

	AddBank(bank *model.Bank) *ErrPixKeyRepository

	AddAccount(account *model.Account) *ErrPixKeyRepository

	FindAccountById(id string) (*model.Account, *ErrPixKeyRepository)
}

type ErrPixKeyRepository struct {
	message string
}

func (err *ErrPixKeyRepository) Error() string {
	return fmt.Sprintf("pix key repository: %s", err.message)
}

var (
	ErrPixKeyAlreadyExists = ErrPixKeyRepository{"pix key already exists"}
	ErrPixKeyNotFound      = ErrPixKeyRepository{"pix key not found"}
	ErrAccountNotFound     = ErrPixKeyRepository{"account not found"}
	ErrOnAddBank           = ErrPixKeyRepository{"cannot add bank"}
	ErrOnAddAccount        = ErrPixKeyRepository{"cannot add account"}
)
