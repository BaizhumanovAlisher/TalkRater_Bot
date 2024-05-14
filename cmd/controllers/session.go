package controllers

import (
	"errors"
	"gorm.io/gorm"
	"talk_rater_bot/internal/data"
)

func (c *Controller) SaveSession(session *data.Session) error {
	return c.db.Save(session).Error
}

func (c *Controller) RetrieveSession(chatID int64) (*data.Session, bool, error) {
	tx := c.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return nil, false, tx.Error
	}

	var session data.Session
	if err := tx.First(&session, "chat_id = ?", chatID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}

	if err := tx.Delete(&session).Error; err != nil {
		tx.Rollback()
		return nil, false, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, false, err
	}

	return &session, true, nil
}
