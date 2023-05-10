package transaction

import (
	"log"

	"gorm.io/gorm"
)

type gormTransactionManager struct {
	db  *gorm.DB
	log log.Logger
}

func NewGormTransactionManager(data *Data, logger log.Logger) TransactionManager {
	return &transactionManager{
		data: data,
		log:  logger,
	}
}

func (s *transactionManager) Run(f func() error) error {
	return s.data.Transaction(func(*gorm.DB) error {
		return f()
	})
}
