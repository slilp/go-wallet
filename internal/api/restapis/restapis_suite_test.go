package restapis_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/slilp/go-wallet/internal/api/restapis"
	"github.com/slilp/go-wallet/internal/api/restapis/api_gen"
	"github.com/slilp/go-wallet/internal/server"
	mock_commands "github.com/slilp/go-wallet/internal/services/commands/mocks"
	mock_queries "github.com/slilp/go-wallet/internal/services/queries/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type RestApisTestSuite struct {
	suite.Suite
	server                 *gin.Engine
	mockRegisterService    *mock_commands.MockRegisterService
	mockTransactionService *mock_commands.MockTransactionService
	mockWalletService      *mock_commands.MockWalletService

	mockListTransactionsService *mock_queries.MockListTransactionsService
	mockListWalletsService      *mock_queries.MockListWalletsService
	mockLoginService            *mock_queries.MockLoginService
}

func (suite *RestApisTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())

	mockListTransactionsService := mock_queries.NewMockListTransactionsService(ctrl)
	mockListWalletsService := mock_queries.NewMockListWalletsService(ctrl)
	mockLoginService := mock_queries.NewMockLoginService(ctrl)
	mockRegisterService := mock_commands.NewMockRegisterService(ctrl)
	mockWalletService := mock_commands.NewMockWalletService(ctrl)
	mockTransactionService := mock_commands.NewMockTransactionService(ctrl)

	r := gin.Default()

	// Add middleware to set user ID for secure routes
	r.Use(func(c *gin.Context) {
		c.Set("USER_ID", "<UserID>")
		c.Next()
	})

	api_gen.RegisterHandlers(r, &restapis.HttpServer{
		App: &server.Application{
			Queries: server.Queries{
				ListWalletsService:      mockListWalletsService,
				ListTransactionsService: mockListTransactionsService,
				LoginService:            mockLoginService,
			},
			Commands: server.Commands{
				RegisterService:    mockRegisterService,
				WalletService:      mockWalletService,
				TransactionService: mockTransactionService,
			},
			Utils: server.Utils{
				Validate: validator.New(),
			},
		},
	})

	suite.mockListTransactionsService = mockListTransactionsService
	suite.mockListWalletsService = mockListWalletsService
	suite.mockLoginService = mockLoginService

	suite.mockRegisterService = mockRegisterService
	suite.mockWalletService = mockWalletService
	suite.mockTransactionService = mockTransactionService

	suite.server = r
}

func TestRestApisTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RestApisTestSuite))
}
