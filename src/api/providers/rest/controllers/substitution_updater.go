package controllers

import (
	"github.com/dattito/purrmannplus-backend/api/providers/rest/models"
	"github.com/dattito/purrmannplus-backend/app/commands"
	"github.com/dattito/purrmannplus-backend/utils/logging"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// Adds an account to the substitution updater
func AddAccountToSubstitutionUpdater(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	accountId := claims["account_id"].(string)
	ok, err := commands.ValidAccountId(accountId)
	if err != nil {
		logging.Errorf("Error validating account id: %s", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Something went wrong",
		})
	}
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "account not found",
		})
	}

	m := models.PostAddAccountToSubstitutionRequest{}
	c.BodyParser(&m)

	if m.Username != "" && m.Password != "" {
		user_err, db_err := commands.AddAccountToSubstitutionUpdaterWithCustomCredentials(accountId, m.Username, m.Password)
		if user_err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"error": user_err.Error(),
			})
		}
		if db_err != nil {
			logging.Errorf("Error while adding account to substitution updater with custom credentials: %v", db_err)
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"error": "Something went wrong",
			})
		}
	} else {
		user_err, db_err := commands.AddAccountToSubstitutionUpdater(accountId)
		if db_err != nil {
			logging.Errorf("Error while adding account to substitution updater: %v", db_err)
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"error": "Something went wrong",
			})
		}
		if user_err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"error": user_err.Error(),
			})
		}
	}
	return c.SendStatus(fiber.StatusCreated)
}

// Removes an account from the substitution updater
func RemoveAccountFromSubstitutionUpdater(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	accountId := claims["account_id"].(string)

	err := commands.RemoveAccountFromSubstitutionUpdater(accountId)
	if err != nil {
		logging.Errorf("Error while removing account from substitution updater: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Something went wrong",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
