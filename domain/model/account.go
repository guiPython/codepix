package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type Account struct {
	Base      `valid:"required"`
	OwnerName string    `json:"owner_name" gorm:"type:varchar(100);notnull" valid:"notnull"`
	BankID    string    `gorm:"column:bank_id;type:uuid;not null" valid:"-"`
	Bank      *Bank     `json:"bank" valid:"-"`
	Number    string    `json:"number" gorm:"type:varchar(20)" valid:"notnull"`
	PixKeys   []*PixKey `gorm:"ForeignKey:AccountID" valid:"-"`
}

func (account *Account) isValid() error {
	if _, err := govalidator.ValidateStruct(account); err != nil {
		return err
	}

	return nil
}

func NewAccount(ownerName, number string, bank *Bank) (*Account, error) {
	account := Account{
		OwnerName: ownerName,
		Number:    number,
		Bank:      bank,
		BankID:    bank.ID,
	}

	account.ID = uuid.NewV4().String()
	account.CreatedAt = time.Now()

	if err := account.isValid(); err != nil {
		return nil, err
	}

	return &account, nil
}
