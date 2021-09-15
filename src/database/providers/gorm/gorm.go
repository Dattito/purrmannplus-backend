package gorm

import (
	"errors"

	"github.com/datti-to/purrmannplus-backend/config"
	db_errors "github.com/datti-to/purrmannplus-backend/database/errors"
	provider_models "github.com/datti-to/purrmannplus-backend/database/models"
	"github.com/datti-to/purrmannplus-backend/database/providers/gorm/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormProvider struct {
	DB *gorm.DB
}

func NewGormProvider() (*GormProvider, error) {

	type Open func(string) gorm.Dialector
	var o Open

	switch config.DATABASE_TYPE {
	case "POSTGRES":
		o = postgres.Open
	case "MYSQL":
		o = mysql.Open
	case "SQLITE":
		o = sqlite.Open
	default:
		return &GormProvider{}, errors.New("DATABASE_TYPE env has to one of ('POSTGRES', 'MYSQL', 'SQLITE')")
	}

	db, err := gorm.Open(o(config.DATABASE_URI), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return &GormProvider{}, err
	}

	return &GormProvider{DB: db}, nil
}

func (g *GormProvider) CreateTables() error {
	var err error
	err = g.DB.AutoMigrate(&models.AccountDB{})
	if err != nil {
		return err
	}

	err = g.DB.AutoMigrate(&models.AccountInfoDB{})
	if err != nil {
		return err
	}

	err = g.DB.AutoMigrate(&models.SubstitutionDB{})
	if err != nil {
		return err
	}
	return nil
}

func (g *GormProvider) CloseDB() error {
	dialect, err := g.DB.DB()
	if err != nil {
		return err
	}
	defer dialect.Close()

	return nil
}

func (g *GormProvider) AddAccount(authId, authPw string) (provider_models.AccountDBModel, error) {

	accdb := models.AccountDB{
		AuthId: authId,
		AuthPw: authPw,
	}
	err := g.DB.Create(&accdb).Error
	return models.AccountDBToAccountDBModel(accdb), err
}

func (g *GormProvider) GetAccount(id string) (provider_models.AccountDBModel, error) {

	accdb := models.AccountDB{}

	err := g.DB.First(&accdb, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return provider_models.AccountDBModel{}, &db_errors.ErrRecordNotFound
		}
		return provider_models.AccountDBModel{}, err
	}

	return models.AccountDBToAccountDBModel(accdb), nil
}

func (g *GormProvider) GetAccountByCredentials(authId, authPw string) (provider_models.AccountDBModel, error) {

	accdb := models.AccountDB{}

	err := g.DB.Where("auth_id = ? AND auth_pw = ?", authId, authPw).First(&accdb).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return provider_models.AccountDBModel{}, &db_errors.ErrRecordNotFound
		}
		return provider_models.AccountDBModel{}, err
	}

	return models.AccountDBToAccountDBModel(accdb), nil
}

func (g *GormProvider) GetAccounts() ([]provider_models.AccountDBModel, error) {

	accdb := []models.AccountDB{}

	err := g.DB.Find(&accdb).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []provider_models.AccountDBModel{}, &db_errors.ErrRecordNotFound
		}
		return []provider_models.AccountDBModel{}, err
	}

	var accs []provider_models.AccountDBModel
	for _, v := range accdb {
		accs = append(accs, models.AccountDBToAccountDBModel(v))
	}

	return accs, nil
}

func (g *GormProvider) DeleteAccount(id string) error {

	accdb := models.AccountDB{}

	err := g.DB.First(&accdb, id).Error

	if err != nil {
		return err
	}

	return g.DB.Delete(&accdb).Error
}

func (g *GormProvider) AddAccountInfo(accountId, phoneNumber string) (provider_models.AccountInfoDBModel, error) {

	accInfo := models.AccountInfoDB{
		AccountId:   accountId,
		PhoneNumber: phoneNumber,
	}
	err := g.DB.Create(&accInfo).Error
	return models.AccountInfoDBToAccountInfoDBModel(accInfo), err
}

func (g *GormProvider) GetAccountInfo(accountId string) (provider_models.AccountInfoDBModel, error) {
	accInfo := models.AccountInfoDB{}
	err := g.DB.Where("account_id = ?", accountId).First(&accInfo).Error
	if err != nil {
		return provider_models.AccountInfoDBModel{}, err
	}
	return models.AccountInfoDBToAccountInfoDBModel(accInfo), err
}

func (g *GormProvider) AddSubstitution(accountId string, entries map[string][]string) (provider_models.SubstitutionDBModel, error) {

	subdb := models.SubstitutionDB{
		AccountId: accountId,
		Entries:   entries,
	}

	err := g.DB.Create(&subdb).Error

	return models.SubstitutionDBToSubstitutionDBModel(subdb), err
}