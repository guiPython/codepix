package usecase

import (
	"errors"

	"github.com/guiPython/codepix/domain/model"
	ports "github.com/guiPython/codepix/domain/ports/repository"
)

type PixKeyService struct {
	repository ports.IPixKeyRepository
}

func NewPixKeyService(repository ports.IPixKeyRepository) *PixKeyService {
	return &PixKeyService{repository: repository}
}

func (srvc *PixKeyService) RegisterKey(rawKind int32, key, accountId string) (*model.PixKey, error) {
	account, repositoryError := srvc.repository.FindAccountById(accountId)
	switch {
	case errors.As(repositoryError, ports.ErrPixKeyRepository{}):
		return nil, repositoryError
	case repositoryError != nil:
		panic(repositoryError)
	}

	pixKey, modelError := model.NewPixKey(key, accountId, rawKind, account)
	switch {
	case errors.As(modelError, model.InvalidPixKey{}):
		return nil, modelError
	case modelError != nil:
		panic(modelError)
	}

	pixKey, repositoryError = srvc.repository.FindPixKeyByKind(key, pixKey.Kind)
	switch {
	case errors.As(repositoryError, ports.ErrPixKeyRepository{}):
		return nil, repositoryError
	case repositoryError != nil:
		panic(repositoryError)
	}

	pixKey, repositoryError = srvc.repository.RegisterPixKey(pixKey)
	switch {
	case errors.As(repositoryError, ports.ErrPixKeyRepository{}):
		return nil, repositoryError
	case repositoryError != nil:
		panic(repositoryError)
	}

	return pixKey, nil
}

func (srvc *PixKeyService) FindKey(rawKind int32, key string) (*model.PixKey, error) {
	var kind model.Kind
	if kind = model.NewKind(rawKind); kind == model.INVALID {
		return nil, &model.ErrInvalidPixKeyKind
	}

	pixKey, err := srvc.repository.FindPixKeyByKind(key, kind)
	switch {
	case errors.As(err, ports.ErrPixKeyRepository{}):
		return nil, err
	case err != nil:
		panic(err)
	}

	return pixKey, nil
}
