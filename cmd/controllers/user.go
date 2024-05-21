package controllers

import (
	"errors"
	"gorm.io/gorm"
	"talk_rater_bot/internal/data"
)

func (c *Controller) UserExists(id int64) (bool, error) {
	user := data.User{ID: id}

	result := c.db.First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}

	return user.IdentityInfo != "", nil
}

func (c *Controller) SaveUser(user *data.User) error {
	return c.db.Save(user).Error
}
