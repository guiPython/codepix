package repository

import (
	"fmt"

	"github.com/guiPython/codepix/domain/model"
)

type ITransactionRepository interface {
	RegisterTransaction(transaction *model.Transaction) (*model.Transaction, *ErrTransactionRepository)

	SaveTransaction(trasaction *model.Transaction) *ErrTransactionRepository

	FindTransactionById(id string) (*model.Transaction, *ErrTransactionRepository)
}

type ErrTransactionRepository struct {
	message string
}

func (err *ErrTransactionRepository) Error() string {
	return fmt.Sprintf("transaction repository: %s", err.message)
}

var (
	ErrTransactionNotFound = ErrTransactionRepository{"transaction not found"}
)
