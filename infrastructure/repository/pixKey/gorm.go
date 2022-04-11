package repository

import (
	"github.com/guiPython/codepix/domain/model"
	port "github.com/guiPython/codepix/domain/ports/repository"
	"github.com/jinzhu/gorm"
)

type PixKeyRepository struct {
	database *gorm.DB
}

func NewPixKeyRepository(database *gorm.DB) port.IPixKeyRepository {
	repository := PixKeyRepository{database: database}
	return &repository
}

func (repository *PixKeyRepository) RegisterPixKey(pixKey *model.PixKey) (*model.PixKey, *port.ErrPixKeyRepository) {
	if err := repository.database.Create(pixKey).Error; err != nil {
		return nil, &port.ErrPixKeyAlreadyExists
	}

	return pixKey, nil
}

func (repository *PixKeyRepository) FindPixKeyByKind(key string, kind model.Kind) (*model.PixKey, *port.ErrPixKeyRepository) {
	var pixKey, nullPixKey model.PixKey
	repository.database.Preload("Account.Bank").First(&pixKey, "kind = ? and key = ?", kind, key)

	if pixKey == nullPixKey {
		return nil, &port.ErrPixKeyNotFound
	}

	return &pixKey, nil
}

func (repository *PixKeyRepository) AddBank(bank *model.Bank) *port.ErrPixKeyRepository {
	err := repository.database.Create(bank).Error
	if err != nil {
		return &port.ErrOnAddBank
	}

	return nil
}

func (repository *PixKeyRepository) AddAccount(account *model.Account) *port.ErrPixKeyRepository {
	if err := repository.database.Create(account).Error; err != nil {
		return &port.ErrOnAddAccount
	}

	return nil
}

func (repository *PixKeyRepository) FindAccountById(id string) (*model.Account, *port.ErrPixKeyRepository) {
	var account model.Account
	repository.database.Preload("Bank").First(&account, "id = ?", id)

	if account.ID == "" {
		return nil, &port.ErrAccountNotFound
	}
	return &account, nil
}
