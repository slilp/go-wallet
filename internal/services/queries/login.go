package queries

import (
	"fmt"
	"log"

	"github.com/slilp/go-wallet/internal/api/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/config"
	"github.com/slilp/go-wallet/internal/repositories"
	"github.com/slilp/go-wallet/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=./login.go -destination=./mocks/mock_login_service.go -package=mock_queries
type LoginService interface {
	Handle(username, password string) (*api_gen.LoginResponseData, error)
}

type loginService struct {
	userRepo repositories.UserRepository
}

func NewLoginService(userRepo repositories.UserRepository) LoginService {
	return &loginService{userRepo: userRepo}
}
func (r *loginService) Handle(email, password string) (*api_gen.LoginResponseData, error) {
	userInfo, err := r.userRepo.QueryByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(password))
	if err != nil {
		log.Printf("Password mismatch for user %s: %v", email, err)
		return nil, err
	}

	accessToken, err := utils.GenerateToken(userInfo.ID, "access", config.Config.AccessTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate access token")
	}

	return &api_gen.LoginResponseData{
		AccessToken: accessToken,
		Email:       userInfo.Email,
		DisplayName: userInfo.DisplayName,
		UserId:      userInfo.ID,
	}, nil
}
