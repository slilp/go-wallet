package repositories_test

// Basic imports
import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/slilp/go-wallet/internal/repositories"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setUpMockDb() (sqlmock.Sqlmock, *gorm.DB) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	return mock, db
}

type UserRepositoryTestSuite struct {
	suite.Suite
	sqlMock  sqlmock.Sqlmock
	userRepo repositories.UserRepository
}

type WalletRepositoryTestSuite struct {
	suite.Suite
	sqlMock    sqlmock.Sqlmock
	walletRepo repositories.WalletRepository
}

type TransactionRepositoryTestSuite struct {
	suite.Suite
	sqlMock         sqlmock.Sqlmock
	transactionRepo repositories.TransactionRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	sqlMock, db := setUpMockDb()
	suite.sqlMock = sqlMock
	suite.userRepo = repositories.NewUserRepository(db)
}

func (suite *WalletRepositoryTestSuite) SetupTest() {
	sqlMock, db := setUpMockDb()
	suite.sqlMock = sqlMock
	suite.walletRepo = repositories.NewWalletRepository(db)
}

func (suite *TransactionRepositoryTestSuite) SetupTest() {
	sqlMock, db := setUpMockDb()
	suite.sqlMock = sqlMock
	suite.transactionRepo = repositories.NewTransactionRepository(db)
}

func TestRepositoryTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UserRepositoryTestSuite))
	suite.Run(t, new(WalletRepositoryTestSuite))
	suite.Run(t, new(TransactionRepositoryTestSuite))
}
