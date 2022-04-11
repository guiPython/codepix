package usecase

import (
	"errors"

	"github.com/guiPython/codepix/domain/model"
	"github.com/guiPython/codepix/domain/ports/repository"
	ports "github.com/guiPython/codepix/domain/ports/repository"
)

type TrasactionService struct {
	transactions ports.ITransactionRepository
	pixKeys      ports.IPixKeyRepository
}

func execute(function func() error) error {
	err := function()
	switch {
	case errors.As(err, &model.ErrInvalidTransaction{}):
		return err
	case errors.As(err, &model.ErrInvalidTransactionOperation{}):
		return err
	case err != nil:
		panic(err)
	}
	return nil
}

func (srvc *TrasactionService) Register(id, description, accountId, pixKeyto string, pixKeyKindTo int32, amount float64) (*model.Transaction, error) {
	account, pixKeyRepositoryError := srvc.pixKeys.FindAccountById(accountId)
	switch {
	case errors.As(pixKeyRepositoryError, &repository.ErrTransactionRepository{}):
		return nil, pixKeyRepositoryError
	case pixKeyRepositoryError != nil:
		panic(pixKeyRepositoryError)
	}

	kind := model.NewKind(pixKeyKindTo)
	if kind == model.INVALID {
		return nil, &model.ErrInvalidPixKeyKind
	}

	pixKey, pixKeyRepositoryError := srvc.pixKeys.FindPixKeyByKind(pixKeyto, kind)
	switch {
	case errors.As(pixKeyRepositoryError, &repository.ErrPixKeyRepository{}):
		return nil, pixKeyRepositoryError
	case pixKeyRepositoryError != nil:
		panic(pixKeyRepositoryError)
	}

	transaction, modelError := model.NewTransaction(account, amount, pixKey, description, "")
	switch {
	case errors.As(modelError, &model.ErrInvalidTransaction{}):
		return nil, modelError
	case modelError != nil:
		panic(modelError)
	}

	transactionRepositoryError := srvc.transactions.SaveTransaction(transaction)
	switch {
	case errors.As(transactionRepositoryError, &repository.ErrTransactionRepository{}):
		return nil, transactionRepositoryError
	case transactionRepositoryError != nil:
		panic(transactionRepositoryError)
	}

	return transaction, nil

}

func (srvc *TrasactionService) Confirm(trasactionId string) (*model.Transaction, error) {
	transaction, repositoryError := srvc.transactions.FindTransactionById(trasactionId)
	switch {
	case errors.As(repositoryError, &repository.ErrTransactionRepository{}):
		return nil, repositoryError
	case repositoryError != nil:
		panic(repositoryError)
	}

	modelError := execute(transaction.Confirm)
	if modelError != nil {
		return nil, modelError
	}

	repositoryError = srvc.transactions.SaveTransaction(transaction)
	switch {
	case errors.As(repositoryError, &repository.ErrTransactionRepository{}):
		return nil, repositoryError
	case repositoryError != nil:
		panic(repositoryError)
	}

	return transaction, nil
}

func (srvc *TrasactionService) Complete(trasactionId string) (*model.Transaction, error) {
	transaction, repositoryError := srvc.transactions.FindTransactionById(trasactionId)
	switch {
	case errors.As(repositoryError, &repository.ErrTransactionRepository{}):
		return nil, repositoryError
	case repositoryError != nil:
		panic(repositoryError)
	}

	execute(transaction.Complete)

	repositoryError = srvc.transactions.SaveTransaction(transaction)
	switch {
	case errors.As(repositoryError, &repository.ErrTransactionRepository{}):
		return nil, repositoryError
	case repositoryError != nil:
		panic(repositoryError)
	}

	return transaction, nil
}

func (srvc *TrasactionService) Error(trasactionId, reason string) (*model.Transaction, error) {
	transaction, repositoryError := srvc.transactions.FindTransactionById(trasactionId)
	switch {
	case errors.As(repositoryError, &repository.ErrTransactionRepository{}):
	case repositoryError != nil:
		panic(repositoryError)
	}

	modelError := transaction.Cancel(reason)
	switch {
	case errors.As(modelError, &model.ErrInvalidTransaction{}):
		return nil, modelError
	case errors.As(modelError, &model.ErrInvalidTransactionOperation{}):
		return nil, modelError
	case modelError != nil:
		panic(modelError)
	}

	repositoryError = srvc.transactions.SaveTransaction(transaction)
	switch {
	case errors.As(repositoryError, &repository.ErrTransactionRepository{}):
	case repositoryError != nil:
		panic(repositoryError)
	}

	return transaction, nil
}
