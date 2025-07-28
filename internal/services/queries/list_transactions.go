package queries

import (
	"github.com/aarondl/null/v9"
	"github.com/slilp/go-wallet/internal/port/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/repositories"
)

//go:generate mockgen -source=./list_transactions.go -destination=./mocks/mock_list_transactions_service.go -package=mock_queries
type ListTransactionsService interface {
	Handle(userId, walletId string, page, limit int) (int64, []api_gen.TransactionResponseData, error)
}

type listTransactionsService struct {
	walletRepo      repositories.WalletRepository
	transactionRepo repositories.TransactionRepository
}

func NewListTransactionsService(walletRepo repositories.WalletRepository, transactionRepo repositories.TransactionRepository) ListTransactionsService {
	return &listTransactionsService{walletRepo: walletRepo, transactionRepo: transactionRepo}
}

func (s *listTransactionsService) Handle(userId, walletId string, page, limit int) (int64, []api_gen.TransactionResponseData, error) {

	_, err := s.walletRepo.QueryByIdAndUser(userId, walletId)
	if err != nil {
		return 0, nil, err
	}

	totalCount, err := s.transactionRepo.CountByWalletId(walletId)
	if err != nil {
		return 0, nil, err
	}

	if totalCount == 0 {
		return totalCount, []api_gen.TransactionResponseData{}, nil
	}

	transactions, err := s.transactionRepo.List(walletId, page, limit)
	if err != nil {

		return 0, nil, err
	}

	result := []api_gen.TransactionResponseData{}
	for _, tx := range transactions {
		result = append(result, api_gen.TransactionResponseData{
			Id:           tx.ID,
			FromWalletId: null.StringFromPtr(tx.From).String,
			ToWalletId:   null.StringFromPtr(tx.To).String,
			Amount:       tx.Amount,
			Type:         api_gen.TransactionResponseDataType(tx.Type),
			CreatedAt:    tx.CreatedAt,
		})
	}

	return totalCount, result, nil
}
