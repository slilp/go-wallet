package commands

import (
	"github.com/slilp/go-wallet/internal/port/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/repositories"
	"github.com/slilp/go-wallet/internal/repositories/entity"
)

//go:generate mockgen -source=./wallet.go -destination=./mocks/mock_wallet_service.go -package=mock_commands
type WalletService interface {
	HandleCreate(userId string, req api_gen.WalletRequest) error
	HandleDelete(walletId string) error
	HandleUpdateInfo(walletId string, req api_gen.WalletRequest) error
}

type walletService struct {
	walletRepo repositories.WalletRepository
}

func NewWalletService(walletRepo repositories.WalletRepository) WalletService {
	return &walletService{walletRepo: walletRepo}
}

func (r *walletService) HandleCreate(userId string, req api_gen.WalletRequest) error {
	return r.walletRepo.Create(entity.Wallet{
		UserID:      userId,
		Name:        req.Name,
		Description: req.Description,
	})
}

func (r *walletService) HandleDelete(walletId string) error {
	return r.walletRepo.Delete(walletId)
}

func (r *walletService) HandleUpdateInfo(walletId string, req api_gen.WalletRequest) error {
	return r.walletRepo.UpdateInfo(walletId, req.Name, req.Description)
}
