package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	provider_models "github.com/datti-to/purrmannplus-backend/database/models"
)

type Entries map[string][]string

func (s *Entries) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &s)
		return nil
	case string:
		json.Unmarshal([]byte(v), &s)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (s *Entries) Value() (driver.Value, error) {
	if s != nil {
		b, err := json.Marshal(s)
		if err != nil {
			return nil, err
		}
		return string(b), nil
	}
	return nil, nil
}

type SubstitutionDB struct {
	Model
	AccountId string    `gorm:"account_id;uniqueIndex"`
	AccountDB AccountDB `gorm:"foreignKey:account_id"`
	Entries   *Entries  `gorm:"entries;default:{}"`
}

func (SubstitutionDB) TableName() string {
	return "substitutions"
}

func (s *SubstitutionDB) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &s)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (s *SubstitutionDB) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func SubstitutionDBToSubstitutionDBModel(s SubstitutionDB) provider_models.SubstitutionDBModel {
	return provider_models.SubstitutionDBModel{
		Id:        s.Id,
		AccountId: s.AccountId,
		Entries:   *s.Entries,
	}
}
