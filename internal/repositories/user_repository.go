package repositories

import (
	"log"

	"github.com/slilp/go-wallet/internal/repositories/entity"
	"gorm.io/gorm"
)

//go:generate mockgen -source=./user_repository.go -destination=./mocks/mock_user_repository.go -package=mock_repositories
type UserRepository interface {
	Create(req entity.User) error
	QueryByEmail(email string) (*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(req entity.User) error {
	if err := r.db.Create(&req).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}
	return nil
}

func (r *userRepository) QueryByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where(&entity.User{Email: email}).Take(&user).Error; err != nil {
		log.Printf("Error querying user by email: %v", err)
		return nil, err
	}
	return &user, nil
}
