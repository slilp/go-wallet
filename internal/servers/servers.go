package servers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	postgres2 "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/slilp/go-wallet/internal/config"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type HttpServer struct {
	Services Services
	Utils    Utils
}

type Services struct {
}

type Utils struct {
	Validate *validator.Validate
}

func NewHttpServer(db *gorm.DB) *HttpServer {
	return &HttpServer{
		Services: Services{},
		Utils: Utils{
			Validate: validator.New(),
		},
	}
}

func InitDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", config.Config.DBUsername, config.Config.DBPassword, config.Config.DBHost, config.Config.DBName, config.Config.DBMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitMigrations(db *gorm.DB) error {
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
