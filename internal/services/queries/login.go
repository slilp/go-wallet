package queries

import (
	"log"

	"github.com/slilp/go-wallet/internal/repositories"
	"github.com/slilp/go-wallet/internal/repositories/entity"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=./login.go -destination=./mocks/mock_login_service.go -package=mock_queries
type LoginService interface {
	Handle(username, password string) (*entity.User, error)
}

type loginService struct {
	userRepo repositories.UserRepository
}

func NewLoginService(userRepo repositories.UserRepository) LoginService {
	return &loginService{userRepo: userRepo}
}
func (r *loginService) Handle(email, password string) (*entity.User, error) {
	userInfo, err := r.userRepo.QueryByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(password))
	if err != nil {
		log.Printf("Password mismatch for user %s: %v", email, err)
		return nil, err
	}

	return userInfo, nil
}
