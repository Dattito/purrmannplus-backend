package commands

import (
	"github.com/datti-to/purrmannplus-backend/app/models"
	"github.com/datti-to/purrmannplus-backend/database"
)

func AddAccountInfo(accountId, phoneNumber string) (models.AccountInfo, error) {
	_, err := models.NewAccountInfo(models.Account{Id: accountId}, phoneNumber)
	if err != nil {
		return models.AccountInfo{}, err
	}

	ai, err := database.DB.AddAccountInfo(accountId, phoneNumber)
	if err != nil {
		return models.AccountInfo{}, err
	}

	return models.AccountInfoDBModelToAccount(ai)
}