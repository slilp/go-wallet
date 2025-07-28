package commands

import (
	"github.com/slilp/go-wallet/internal/repositories"
)

//go:generate mockgen -source=./transaction.go -destination=./mocks/mock_transaction_service.go -package=mock_commands
type TransactionService interface {
	HandleTransferBalance(userId, from, to string, amount float64) error
	HandleDepositWithDrawBalance(userId, walletId string, amount float64) error
}

type transactionService struct {
	transactionRepo repositories.TransactionRepository
}

func NewTransactionService(transactionRepo repositories.TransactionRepository) TransactionService {
	return &transactionService{transactionRepo: transactionRepo}
}

func (r *transactionService) HandleTransferBalance(userId, from, to string, amount float64) error {
	return r.transactionRepo.UpdateTransferTransaction(userId, from, to, amount)
}

func (r *transactionService) HandleDepositWithDrawBalance(userId, walletId string, amount float64) error {
	return r.transactionRepo.UpdateBalanceTransaction(userId, walletId, amount)
}
