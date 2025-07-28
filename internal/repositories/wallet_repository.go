package repositories

import (
	"log"

	"github.com/slilp/go-wallet/internal/repositories/entity"
	"gorm.io/gorm"
)

//go:generate mockgen -source=./wallet_repository.go -destination=./mocks/mock_wallet_repository.go -package=mock_repositories
type WalletRepository interface {
	Create(req entity.Wallet) error
	UpdateInfo(id, name string, desc *string) error
	Delete(walletId string) error
	ListAll(userId string) ([]entity.Wallet, error)
	QueryByIdAndUser(userId, walletId string) (*entity.Wallet, error)
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) Create(req entity.Wallet) error {
	if err := r.db.Create(&req).Error; err != nil {
		log.Printf("Create error: %v", err)
		return err
	}
	return nil
}

func (r *walletRepository) UpdateInfo(id, name string, desc *string) error {
	if err := r.db.Model(&entity.Wallet{}).
		Where(&entity.Wallet{ID: id}).
		Updates(entity.Wallet{Name: name, Description: desc}).Error; err != nil {
		log.Printf("UpdateInfo error: %v", err)
		return err
	}
	return nil
}

func (r *walletRepository) Delete(walletId string) error {
	wallet := entity.Wallet{ID: walletId}
	if err := r.db.Delete(&wallet).Error; err != nil {
		log.Printf("Delete error: %v", err)
		return err
	}
	return nil
}

func (r *walletRepository) ListAll(userId string) ([]entity.Wallet, error) {
	var wallets []entity.Wallet
	if err := r.db.Where(&entity.Wallet{UserID: userId}).Find(&wallets).Error; err != nil {
		log.Printf("ListAll error: %v", err)
		return nil, err
	}
	return wallets, nil
}

func (r *walletRepository) QueryByIdAndUser(userId, walletId string) (*entity.Wallet, error) {
	var wallet entity.Wallet
	if err := r.db.Where(&entity.Wallet{UserID: userId, ID: walletId}).First(&wallet).Error; err != nil {
		log.Printf("QueryByIdAndUser error: %v", err)
		return nil, err
	}
	return &wallet, nil
}
