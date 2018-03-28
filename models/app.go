package models

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// App ...
type App struct {
	gorm.Model
	Plan     string `json:"plan" gorm:"type:varchar(100)"`
	AppSlug  string `json:"app_slug" gorm:"type:varchar(100)"`
	APIToken string `json:"-" gorm:"type:char(26)"`
}

// Exists ...
func (a *App) Exists() (bool, error) {
	if err := DB.Find(a, a).Error; err != nil {
		if errors.Cause(err) == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
