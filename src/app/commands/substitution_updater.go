package commands

import (
	"errors"
	"fmt"
	"log"

	"github.com/datti-to/purrmannplus-backend/app/models"
	"github.com/datti-to/purrmannplus-backend/config"
	"github.com/datti-to/purrmannplus-backend/database"
	db_errors "github.com/datti-to/purrmannplus-backend/database/errors"
	"github.com/datti-to/purrmannplus-backend/services/hpg"
	"github.com/datti-to/purrmannplus-backend/services/scheduler"
	"github.com/datti-to/purrmannplus-backend/services/signal_message_sender"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func differenceAmount(newSubstituations map[string][]string, oldSubstituations map[string][]string) map[string][]string {
	s := map[string][]string{}

	for day, lessons := range newSubstituations {
		if len(oldSubstituations[day]) <= 0 {
			s[day] = lessons
			continue
		} else {
			for _, lesson := range lessons {
				if !contains(oldSubstituations[day], lesson) {
					s[day] = append(s[day], lesson)
				}
			}
		}
	}
	return s
}

func substituationToTextMessage(substitution map[string][]string) string {
	if len(substitution) == 0 {
		return "Du hast keine neuen Vertretungen"
	}

	var text string = "Du hast neue Vertretungen: \n"

	for day, lessons := range substitution {
		text += fmt.Sprintf("\n%s:\n", day)
		for _, lesson := range lessons {
			text += fmt.Sprintf("%s\n", lesson)
		}
	}
	return text
}

func AddToSubstitutionUpdater(accountId string) error {
	ai, err := database.DB.GetAccountInfo(accountId)
	if err != nil {
		if errors.Is(err, &db_errors.ErrRecordNotFound) {
			return errors.New("phone number has to be added first")
		}
		return err
	}

	if ai.PhoneNumber == "" {
		return errors.New("phone number has to be added first")
	}

	if _, err := database.DB.GetSubstitutions(accountId); err == nil || !errors.Is(err, &db_errors.ErrRecordNotFound) {
		return errors.New("substitutions already exist")
	}

	_, err = database.DB.SetSubstitutions(accountId, map[string][]string{}, true)
	if err != nil {
		return err
	}

	return UpdateSubstitutionsByAccountId(accountId)
}

func RemoveFromSubstitutionUpdater(accountId string) error {
	return database.DB.RemoveAccountFromSubstitutionUpdater(accountId)
}

// Updates the substitutions for a given account and it's relevant data and sends a message via signal
func UpdateSubstitutions(m models.SubstitutionUpdateInfos) error {
	log.Printf("Updating substitutions of account %s (id: %s)", m.AuthId, m.AccountId)
	mayNewSubstitutions, err := hpg.GetSubstituationOfStudent(m.AuthId, m.AuthPw)
	if err != nil {
		return err
	}

	old_substitutions := m.Entries

	newSubstitutions := differenceAmount(mayNewSubstitutions, old_substitutions)

	// If there are no new substitutions, we don't need to do anything
	if len(newSubstitutions) == 0 {
		return nil
	}

	_, err = database.DB.SetSubstitutions(m.AccountId, mayNewSubstitutions, false)
	if err != nil {
		return err
	}
	// Send a message to the user if there are new substitutions
	if m.NotSetYet {
		return nil
	}

	if err := signal_message_sender.SignalMessageSender.Send(substituationToTextMessage(newSubstitutions), m.PhoneNumber); err != nil {
		return err
	}

	return nil
}

func UpdateSubstitutionsByAccountId(accountId string) error {
	mdb, err := database.DB.GetAccountCredentialsAndPhoneNumberAndSubstitutions(accountId)
	if err != nil {
		return err
	}

	m := models.AccountCredentialsAndPhoneNumberAndSubstitutionsDBModelToSubstitutionUpdateInfos(&mdb)

	return UpdateSubstitutions(m)
}

func UpdateAllSubstitutions() error {
	mdbs, err := database.DB.GetAllAccountCredentialsAndPhoneNumberAndSubstitutions()
	if err != nil {
		return err
	}

	errCount := 0

	for _, mdb := range mdbs {
		m := models.AccountCredentialsAndPhoneNumberAndSubstitutionsDBModelToSubstitutionUpdateInfos(&mdb)
		err := UpdateSubstitutions(m)
		if err != nil {
			log.Printf("Error updating substitutions for account %s: %s", mdb.AccountId, err.Error())
			errCount++
			if errCount > config.MAX_ERROS_TO_STOP_UPDATING_SUBSTITUTIONS {
				return errors.New("got too many errors updating substitutions, stopping")
			}
		}
	}
	return nil
}

func EnableSubstitutionUpdater() {
	scheduler.AddJob(config.SUBSTITUTIONS_UPDATECRON, func() {
		if err := UpdateAllSubstitutions(); err != nil {
			log.Println(err.Error())
		}
	})
}
