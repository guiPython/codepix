package model

import (
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type Kind byte
type Status byte

const (
	EMAIL    Kind   = iota
	CPF      Kind   = 1
	INVALID  Kind   = 2
	ACTIVE   Status = iota
	INACTIVE Status = 1
)

func NewKind(raw int32) Kind {
	switch raw {
	case 0:
		return EMAIL
	case 1:
		return CPF
	default:
		return INVALID
	}
}

type PixKey struct {
	Base      `valid:"required"`
	Kind      Kind     `json:"kind" gorm:"type:smallint" valid:"notnull"`
	Key       string   `json:"key" gorm:"type:varchar(100)" valid:"notnull"`
	AccountID string   `gorm:"column:account_id;type:uuid;not null" valid:"-"`
	Account   *Account `valid:"-"`
	Status    Status   `json:"status" gorm:"type:smallint" valid:"notnull"`
}

type InvalidPixKey struct{ message string }

func (err *InvalidPixKey) Error() string {
	return fmt.Sprintf("pix key: invalid pix key: %s", err.message)
}

type InvalidChangePixKeyStatus struct{ message string }

func (err *InvalidChangePixKeyStatus) Error() string {
	return fmt.Sprintf("pix key: invalid change status: %s", err.message)
}

var (
	ErrInvalidPixKeyStatus = InvalidPixKey{"invalid status must be ACTIVE | INACTIVE"}
	ErrInvalidPixKeyKind   = InvalidPixKey{"invalid kind must be EMAIL | CPF"}

	ErrStatusAlreadyActive   = InvalidChangePixKeyStatus{"status already active"}
	ErrStatusAlreadyInactive = InvalidChangePixKeyStatus{"pix key: status already inactive"}
)

func (pixKey *PixKey) isValid() *InvalidPixKey {
	if _, err := govalidator.ValidateStruct(pixKey); err != nil {
		return &InvalidPixKey{err.Error()}
	}

	if pixKey.Kind == INVALID {
		return &ErrInvalidPixKeyKind
	}

	if pixKey.Status != ACTIVE && pixKey.Status != INACTIVE {
		return &ErrInvalidPixKeyStatus
	}

	return nil
}

func NewPixKey(key, accountID string, kind int32, account *Account) (*PixKey, *InvalidPixKey) {
	pixKey := PixKey{
		Key:       key,
		Kind:      NewKind(kind),
		AccountID: accountID,
		Account:   account,
	}

	pixKey.ID = uuid.NewV4().String()
	pixKey.CreatedAt = time.Now()

	if err := pixKey.isValid(); err != nil {
		return nil, err
	}

	return &pixKey, nil
}

func (pixKey *PixKey) Active() error {
	if pixKey.Status == INACTIVE {
		pixKey.Status = ACTIVE
		pixKey.UpdatedAt = time.Now()
		if err := pixKey.isValid(); err != nil {
			return err
		}
		return nil
	}
	return &ErrStatusAlreadyActive
}

func (pixKey *PixKey) Inactive() error {
	if pixKey.Status == ACTIVE {
		pixKey.Status = INACTIVE
		pixKey.UpdatedAt = time.Now()
		if err := pixKey.isValid(); err != nil {
			return err
		}
		return nil
	}
	return &ErrStatusAlreadyInactive
}
