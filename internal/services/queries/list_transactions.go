package queries

import (
	"github.com/slilp/go-wallet/internal/adapters"
	"github.com/slilp/go-wallet/internal/port/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/repositories"
)

//go:generate mockgen -source=./list_transactions.go -destination=./mocks/mock_list_transactions_service.go -package=mock_queries
type ListTransactionsService interface {
	Handle(userId, walletId string, page, limit int) ([]api_gen.TransactionResponseData, error)
}

type listTransactionsService struct {
	transactionRepo repositories.TransactionRepository
	redisAdapter    adapters.RedisAdapter
}

func NewListTransactionsService(transactionRepo repositories.TransactionRepository, redisAdapter adapters.RedisAdapter) ListTransactionsService {
	return &listTransactionsService{transactionRepo: transactionRepo, redisAdapter: redisAdapter}
}

func (s *listTransactionsService) Handle(userId, walletId string, page, limit int) ([]api_gen.TransactionResponseData, error) {
	return nil, nil
}
