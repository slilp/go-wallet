package commands

import (
	"github.com/slilp/go-wallet/internal/repositories"
)

//go:generate mockgen -source=./wallet.go -destination=./mocks/mock_wallet_service.go -package=mock_commands
type TransactionService interface {
	HandleTransferBalance(from, to string, amount float64) error
	HandleDepositWithDrawBalance(walletId string, amount float64) error
}

type transactionService struct {
	transactionRepo repositories.TransactionRepository
}

func NewTransactionService(transactionRepo repositories.TransactionRepository) TransactionService {
	return &transactionService{transactionRepo: transactionRepo}
}

func (r *transactionService) HandleTransferBalance(from, to string, amount float64) error {
	return r.transactionRepo.UpdateTransferTransaction(from, to, amount)
}

func (r *transactionService) HandleDepositWithDrawBalance(walletId string, amount float64) error {
	return r.transactionRepo.UpdateBalanceTransaction(walletId, amount)
}
