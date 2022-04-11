package repository

import (
	"github.com/guiPython/codepix/domain/model"
	port "github.com/guiPython/codepix/domain/ports/repository"
)

type PixKeyRepositoryInMemory struct {
	banks    []model.Bank
	accounts []model.Account
	keys     []model.PixKey
}

func NewPixKeyRepositoryInMemory() port.IPixKeyRepository {
	banks := []model.Bank{}
	accounts := []model.Account{}
	keys := []model.PixKey{}
	repository := PixKeyRepositoryInMemory{banks: banks, accounts: accounts, keys: keys}
	return &repository
}

func (repository *PixKeyRepositoryInMemory) RegisterPixKey(pixKey *model.PixKey) (*model.PixKey, *port.ErrPixKeyRepository) {
	repository.keys = append(repository.keys, *pixKey)
	return pixKey, &port.ErrPixKeyAlreadyExists
}

func (repository *PixKeyRepositoryInMemory) FindPixKeyByKind(key string, kind model.Kind) (*model.PixKey, *port.ErrPixKeyRepository) {
	for _, pixKey := range repository.keys {
		if pixKey.Kind == kind && pixKey.Key == key {
			return &pixKey, nil
		}
	}
	return nil, &port.ErrPixKeyNotFound
}

func (repository *PixKeyRepositoryInMemory) AddBank(bank *model.Bank) *port.ErrPixKeyRepository {
	repository.banks = append(repository.banks, *bank)
	return nil
}

func (repository *PixKeyRepositoryInMemory) AddAccount(account *model.Account) *port.ErrPixKeyRepository {
	repository.accounts = append(repository.accounts, *account)
	return nil
}

func (repository *PixKeyRepositoryInMemory) FindAccountById(id string) (*model.Account, *port.ErrPixKeyRepository) {
	for _, pixKey := range repository.keys {
		if pixKey.AccountID == id {
			return pixKey.Account, nil
		}
	}
	return nil, &port.ErrAccountNotFound
}
