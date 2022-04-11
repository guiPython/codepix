package model

import (
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type TransactionStatus int

const (
	PENDING              TransactionStatus = iota
	CONFIRMED            TransactionStatus = 1
	COMPLETED            TransactionStatus = 2
	CANCELED             TransactionStatus = 3
	minTransactionAmount float64           = 0
)

type Transaction struct {
	Base              `valid:"required"`
	Amount            float64           `json:"amount" gorm:"type:float" valid:"notnull"`
	Status            TransactionStatus `json:"status" gorm:"type:smallint" valid:"notnull"`
	AccountFromID     string            `gorm:"column:account_from_id;type:uuid;" valid:"notnull"`
	AccountFrom       *Account          `json:"account_from" valid:"notnull"`
	PixKeyIdTo        string            `gorm:"column:pix_key_id_to;type:uuid;" valid:"notnull"`
	PixKeyTo          *PixKey           `json:"pixKey_to" valid:"notnull"`
	Description       string            `json:"description" valid:"notnull"`
	CancelDescription string            `json:"cancel_description" gorm:"type:varchar(255)" valid:"-"`
}

type ErrInvalidTransaction struct{ message string }

func (err *ErrInvalidTransaction) Error() string {
	return fmt.Sprintf("transaction: invalid transaction: %s", err.message)
}

type ErrInvalidTransactionOperation struct{ message string }

func (err *ErrInvalidTransactionOperation) Error() string {
	return fmt.Sprintf("transaction operation: invalid transaction operation: %s", err.message)
}

var (
	errInvalidTransactionStatus  = ErrInvalidTransaction{"status must be PENDING | CONFIRMED | COMPLETED | CANCELED"}
	errInvalidTransactionAmount  = ErrInvalidTransaction{"amount must be bigger 0"}
	errInvalidTransactionAccount = ErrInvalidTransaction{"destiny account equals sender account"}

	errOnCompleteTransaction = ErrInvalidTransactionOperation{"cannot complete canceled or confirmed trasaction"}
	errOnCancelTransaction   = ErrInvalidTransactionOperation{"cannot cancel transaction after she's completed or confirmed"}
	errOnConfirmTransaction  = ErrInvalidTransactionOperation{"cannot confirm transaction before she's pending or canceled"}
)

func (transaction *Transaction) isValid() *ErrInvalidTransaction {
	if _, err := govalidator.ValidateStruct(transaction); err != nil {
		return &ErrInvalidTransaction{err.Error()}
	}

	if transaction.Status != COMPLETED && transaction.Status != CANCELED &&
		transaction.Status != PENDING && transaction.Status != CONFIRMED {
		return &errInvalidTransactionStatus
	}

	if transaction.PixKeyTo.AccountID == transaction.AccountFrom.ID {
		return &errInvalidTransactionAccount
	}

	if transaction.Amount <= minTransactionAmount {
		return &errInvalidTransactionAmount
	}

	return nil
}

func NewTransaction(accountFrom *Account, amount float64, pixKeyTo *PixKey, description, id string) (*Transaction, error) {
	transaction := Transaction{
		AccountFrom:   accountFrom,
		AccountFromID: accountFrom.ID,
		Amount:        amount,
		PixKeyTo:      pixKeyTo,
		PixKeyIdTo:    pixKeyTo.ID,
		Status:        PENDING,
		Description:   description,
	}
	if id == "" {
		transaction.ID = uuid.NewV4().String()
	} else {
		transaction.ID = id
	}
	transaction.CreatedAt = time.Now()

	if err := transaction.isValid(); err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (transaction *Transaction) Complete() error {
	if transaction.Status == PENDING {
		transaction.Status = COMPLETED
		transaction.UpdatedAt = time.Now()
		if err := transaction.isValid(); err != nil {
			return err
		}
		return nil
	}
	return &errOnCompleteTransaction
}

func (transaction *Transaction) Confirm() error {
	if transaction.Status == COMPLETED {
		transaction.Status = CONFIRMED
		transaction.UpdatedAt = time.Now()
		if err := transaction.isValid(); err != nil {
			return err
		}
		return nil
	}
	return &errOnConfirmTransaction
}

func (transaction *Transaction) Cancel(description string) error {
	if transaction.Status == PENDING {
		transaction.Status = CANCELED
		transaction.UpdatedAt = time.Now()
		if err := transaction.isValid(); err != nil {
			return err
		}
		return nil
	}
	return &errOnCancelTransaction
}
