package queries_test

import (
	"testing"

	mock_repositories "github.com/slilp/go-wallet/internal/repositories/mocks"
	"github.com/slilp/go-wallet/internal/services/queries"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type QueriesTestSuite struct {
	suite.Suite
	loginService            queries.LoginService
	listWalletsService      queries.ListWalletsService
	listTransactionsService queries.ListTransactionsService

	mockUserRepo        *mock_repositories.MockUserRepository
	mockWalletRepo      *mock_repositories.MockWalletRepository
	mockTransactionRepo *mock_repositories.MockTransactionRepository
}

func (suite *QueriesTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())

	mockUserRepo := mock_repositories.NewMockUserRepository(ctrl)
	mockWalletRepo := mock_repositories.NewMockWalletRepository(ctrl)
	mockTransactionRepo := mock_repositories.NewMockTransactionRepository(ctrl)
	suite.mockUserRepo = mockUserRepo
	suite.mockWalletRepo = mockWalletRepo
	suite.mockTransactionRepo = mockTransactionRepo

	suite.loginService = queries.NewLoginService(mockUserRepo)
	suite.listWalletsService = queries.NewListWalletsService(mockWalletRepo)
	suite.listTransactionsService = queries.NewListTransactionsService(mockWalletRepo, mockTransactionRepo)
}

func TestQueriesTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(QueriesTestSuite))
}
