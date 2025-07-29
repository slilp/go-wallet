package commands

import (
	"github.com/slilp/go-wallet/internal/api/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/repositories"
	"github.com/slilp/go-wallet/internal/repositories/entity"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=./register.go -destination=./mocks/mock_register_service.go -package=mock_commands
type RegisterService interface {
	Handle(req api_gen.RegisterRequest) error
}

type registerService struct {
	userRepo repositories.UserRepository
}

func NewRegisterService(userRepo repositories.UserRepository) RegisterService {
	return &registerService{userRepo: userRepo}
}

func (r *registerService) Handle(req api_gen.RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return r.userRepo.Create(entity.User{
		Email:       req.Email,
		Password:    string(hashedPassword),
		DisplayName: req.DisplayName,
	})
}
