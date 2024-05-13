package controllers

import (
	"errors"
	"gorm.io/gorm"
	"talk_rater_bot/internal/data"
)

func (c *Controller) UserExists(id int64) (bool, error) {
	var count int64
	result := c.db.Model(&data.User{}).Where("id = ?", id).Count(&count)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, result.Error
	}

	return count > 0, nil
}
