package commands_test

import (
	"testing"

	mock_repositories "github.com/slilp/go-wallet/internal/repositories/mocks"
	"github.com/slilp/go-wallet/internal/services/commands"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type CommandsTestSuite struct {
	suite.Suite
	registerService     commands.RegisterService
	walletService       commands.WalletService
	transactionService  commands.TransactionService
	mockWalletRepo      *mock_repositories.MockWalletRepository
	mockUserRepo        *mock_repositories.MockUserRepository
	mockTransactionRepo *mock_repositories.MockTransactionRepository
}

func (suite *CommandsTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())

	mockUserRepo := mock_repositories.NewMockUserRepository(ctrl)
	mockWalletRepo := mock_repositories.NewMockWalletRepository(ctrl)
	mockTransactionRepo := mock_repositories.NewMockTransactionRepository(ctrl)
	suite.mockUserRepo = mockUserRepo
	suite.mockWalletRepo = mockWalletRepo
	suite.mockTransactionRepo = mockTransactionRepo

	suite.registerService = commands.NewRegisterService(mockUserRepo)
	suite.walletService = commands.NewWalletService(mockWalletRepo)
	suite.transactionService = commands.NewTransactionService(mockTransactionRepo)
}

func TestCommandsTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(CommandsTestSuite))
}
