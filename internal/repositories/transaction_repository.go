package repositories

import (
	"log"

	"github.com/slilp/go-wallet/internal/repositories/entity"
	"gorm.io/gorm"
)

//go:generate mockgen -source=./transaction_repository.go -destination=./mocks/mock_transaction_repository.go -package=mock_repositories
type TransactionRepository interface {
	UpdateBalanceTransaction(walletId string, amount float64) error
	UpdateTransferTransaction(from, to string, amount float64) error
	List(walletId string, page, limit int) ([]entity.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(req entity.Transaction) error {
	if err := r.db.Create(&req).Error; err != nil {
		log.Printf("Create transaction error: %v", err)
		return err
	}

	return nil
}

func (r *transactionRepository) UpdateTransferTransaction(from, to string, amount float64) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Model(&entity.Wallet{}).
			Where(&entity.Wallet{ID: from}).
			UpdateColumn("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			log.Printf("UpdateTransferBalance (from) error: %v", err)
			return err
		}

		if err := tx.Model(&entity.Wallet{}).
			Where(&entity.Wallet{ID: to}).
			UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			log.Printf("UpdateTransferBalance (to) error: %v", err)
			return err
		}

		txRecord := entity.Transaction{
			From:   from,
			To:     to,
			Amount: amount,
			Type:   "transfer",
		}
		if err := tx.Create(&txRecord).Error; err != nil {
			log.Printf("Create transfer transaction error: %v", err)
			return err
		}
		return nil
	}); err != nil {
		log.Printf("UpdateTransferBalance transaction error: %v", err)
		return err
	}
	return nil
}

func (r *transactionRepository) UpdateBalanceTransaction(walletId string, amount float64) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Model(&entity.Wallet{}).
			Where(&entity.Wallet{ID: walletId}).
			UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			log.Printf("UpdateBalance error: %v", err)
			return err
		}

		txRecord := entity.Transaction{
			To:     walletId,
			Amount: amount,
			Type:   "deposit",
		}

		if amount < 0 {
			txRecord = entity.Transaction{
				From:   walletId,
				Amount: amount,
				Type:   "withdraw",
			}
		}

		if err := tx.Create(&txRecord).Error; err != nil {
			log.Printf("Create %s transaction error: %v", txRecord.Type, err)
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (r *transactionRepository) List(walletId string, page, limit int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	offset := (page - 1) * limit

	if err := r.db.Where("from = ? OR to = ?", walletId, walletId).
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		log.Printf("List transactions error: %v", err)
		return nil, err
	}

	return transactions, nil
}
