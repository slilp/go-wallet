package queries

import (
	"github.com/slilp/go-wallet/internal/port/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/repositories"
	"github.com/slilp/go-wallet/internal/repositories/entity"
)

//go:generate mockgen -source=./list_wallets.go -destination=./mocks/mock_list_wallets_service.go -package=mock_queries
type ListWalletsService interface {
	Handle(userId string) ([]api_gen.WalletResponseData, error)
}

type listWalletsService struct {
	walletRepo repositories.WalletRepository
}

func NewListWalletsService(walletRepo repositories.WalletRepository) ListWalletsService {
	return &listWalletsService{walletRepo: walletRepo}
}

func (r *listWalletsService) Handle(userId string) ([]api_gen.WalletResponseData, error) {
	wallets, err := r.walletRepo.ListAll(userId)
	if err != nil {
		return nil, err
	}

	return mapRepoToResponse(wallets), nil
}

func mapRepoToResponse(wallets []entity.Wallet) []api_gen.WalletResponseData {
	response := []api_gen.WalletResponseData{}
	for _, wallet := range wallets {
		response = append(response, api_gen.WalletResponseData{
			Id:          wallet.ID,
			Balance:     wallet.Balance,
			Name:        wallet.Name,
			Description: wallet.Description,
			UpdatedAt:   wallet.UpdatedAt,
		})
	}
	return response
}
