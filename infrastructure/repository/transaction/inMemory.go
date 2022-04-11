package repository

import (
	"github.com/guiPython/codepix/domain/model"
	port "github.com/guiPython/codepix/domain/ports/repository"
)

type TransactionsRepositoryInMemory struct {
	transactions []model.Transaction
}

func NewTransactionRepositoryInMemory() port.ITransactionRepository {
	transactions := []model.Transaction{}
	repository := TransactionsRepositoryInMemory{transactions: transactions}
	return &repository
}

func (repository *TransactionsRepositoryInMemory) RegisterTransaction(transaction *model.Transaction) (*model.Transaction, *port.ErrTransactionRepository) {
	repository.transactions = append(repository.transactions, *transaction)
	return transaction, nil
}

func (repository *TransactionsRepositoryInMemory) SaveTransaction(trasaction *model.Transaction) *port.ErrTransactionRepository {
	return nil
}

func (repository *TransactionsRepositoryInMemory) FindTransactionById(id string) (*model.Transaction, *port.ErrTransactionRepository) {
	for _, transaction := range repository.transactions {
		if transaction.ID == id {
			return &transaction, nil
		}
	}
	return nil, &port.ErrTransactionNotFound
}
