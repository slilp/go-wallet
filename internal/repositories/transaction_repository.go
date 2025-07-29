package repositories

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/aarondl/null/v9"
	"github.com/slilp/go-wallet/internal/consts"
	"github.com/slilp/go-wallet/internal/repositories/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockgen -source=./transaction_repository.go -destination=./mocks/mock_transaction_repository.go -package=mock_repositories
type TransactionRepository interface {
	UpdateBalanceTransaction(userId, walletId string, amount float64) error
	UpdateTransferTransaction(userId, from, to string, amount float64) error
	List(walletId string, page, limit int) ([]entity.Transaction, error)
	CountByWalletId(walletId string) (int64, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) UpdateTransferTransaction(userId, from, to string, amount float64) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {

		var fromWallet entity.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(&entity.Wallet{ID: from, UserID: userId}).
			First(&fromWallet).Error; err != nil {
			log.Printf("Failed to lock (from) wallet: %v", err)
			return err
		}

		if fromWallet.Balance < amount {
			log.Printf("Insufficient balance: wallet %s has %.2f, attempted %.2f", from, fromWallet.Balance, amount)
			return consts.ErrInsufficientBalance
		}

		var toWallet entity.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(&entity.Wallet{ID: to}).
			First(&toWallet).Error; err != nil {
			log.Printf("Failed to lock (to) wallet: %v", err)
			return err
		}

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
			ID:     generateTransactionId(),
			From:   null.StringFrom(from).Ptr(),
			To:     null.StringFrom(to).Ptr(),
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

func (r *transactionRepository) UpdateBalanceTransaction(userId, walletId string, amount float64) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var lockWallet entity.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(&entity.Wallet{ID: walletId, UserID: userId}).
			First(&lockWallet).Error; err != nil {
			log.Printf("Failed to lock wallet: %v", err)
			return err
		}

		txRecord := entity.Transaction{
			ID:     generateTransactionId(),
			To:     null.StringFrom(walletId).Ptr(),
			Amount: amount,
			Type:   "deposit",
		}

		if amount < 0 {
			if lockWallet.Balance < -amount {
				log.Printf("Insufficient balance: wallet %s has %.2f, attempted %.2f", walletId, lockWallet.Balance, amount)
				return consts.ErrInsufficientBalance
			}

			txRecord.To = nil
			txRecord.From = null.StringFrom(walletId).Ptr()
			txRecord.Type = "withdraw"
		}

		if err := tx.Model(&entity.Wallet{}).
			Where(&entity.Wallet{ID: walletId}).
			UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			log.Printf("UpdateBalance error: %v", err)
			return err
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

	if err := r.db.Where(`"from" = ? OR "to" = ?`, walletId, walletId).
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		log.Printf("List transactions error: %v", err)
		return nil, err
	}

	return transactions, nil
}

func (r *transactionRepository) CountByWalletId(walletId string) (int64, error) {
	var count int64
	if err := r.db.Model(&entity.Transaction{}).
		Where(`"from" = ? OR "to" = ?`, walletId, walletId).
		Count(&count).Error; err != nil {
		log.Printf("Count transactions error: %v", err)
		return 0, err
	}
	return count, nil
}

func generateTransactionId() string {
	unixTime := time.Now().Unix()

	timePart := fmt.Sprintf("%010d", unixTime)

	rand.Seed(time.Now().UnixNano())
	randomPart := fmt.Sprintf("%07d", rand.Intn(10000000))

	transactionID := fmt.Sprintf("TRN%s%s", timePart, randomPart)

	return transactionID
}
