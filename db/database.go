package db

import (
	"bank-grpc/model"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectAndMigrate() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", viper.GetString("DB_HOST"), viper.GetString("DB_USER"), viper.GetString("DB_PASSWORD"), viper.GetString("DB_NAME"), viper.GetString("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("unable to connect to postgres: %v", err)
	} else {
		log.Info("Connected to database successfully")
	}
	db.AutoMigrate(&model.Account{})
	db.AutoMigrate(&model.DepositWithdraw{})
	db.AutoMigrate(&model.Transaction{})
	return db, nil
}
