package server

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	postgres2 "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/slilp/go-wallet/internal/config"
	"github.com/slilp/go-wallet/internal/repositories"
	"github.com/slilp/go-wallet/internal/services/commands"
	"github.com/slilp/go-wallet/internal/services/queries"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Application struct {
	Queries  Queries
	Commands Commands
	Utils    Utils
}

type Queries struct {
	ListWalletsService      queries.ListWalletsService
	ListTransactionsService queries.ListTransactionsService
	LoginService            queries.LoginService
}

type Commands struct {
	RegisterService    commands.RegisterService
	WalletService      commands.WalletService
	TransactionService commands.TransactionService
}

type Utils struct {
	Validate *validator.Validate
}

func NewApplicationServer() *Application {

	db, err := initDatabase()
	if err != nil {
		log.Panic(err)
	}

	if err := initMigrations(db); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No new migrations to apply.")
		} else {
			log.Panic("Error applying migrations:", err)
		}
	}

	userRepo := repositories.NewUserRepository(db)
	walletRepo := repositories.NewWalletRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	return &Application{
		Queries: Queries{
			ListWalletsService:      queries.NewListWalletsService(walletRepo),
			ListTransactionsService: queries.NewListTransactionsService(walletRepo, transactionRepo),
			LoginService:            queries.NewLoginService(userRepo),
		},
		Commands: Commands{
			RegisterService: commands.NewRegisterService(userRepo),
			WalletService:   commands.NewWalletService(walletRepo),
		},
		Utils: Utils{
			Validate: validator.New(),
		},
	}
}

func initDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", config.Config.DBUsername, config.Config.DBPassword, config.Config.DBHost, config.Config.DBName, config.Config.DBMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initMigrations(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	driver, err := postgres2.WithInstance(sqlDB, &postgres2.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations",
		"postgres2", driver)
	if err != nil {
		return err
	}

	return m.Up()
}
