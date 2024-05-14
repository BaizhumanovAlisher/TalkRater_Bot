package controllers

import (
	"errors"
	"gorm.io/gorm"
	"talk_rater_bot/internal/data"
)

func (c *Controller) SaveEvaluation(eval *data.Evaluation) error {
	var evaluation data.Evaluation
	result := c.db.Select("id").Where("user_id = ? AND lecture_id = ?", eval.UserID, eval.LectureID).First(&evaluation)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.db.Create(eval).Error
	}

	if result.Error != nil {
		return result.Error
	}

	eval.ID = evaluation.ID
	return c.db.Save(eval).Error
}

func (c *Controller) SaveComment(id int64, comment string) error {
	return c.db.Model(&data.Evaluation{}).Where("ID = ?", id).Update("Comment", comment).Error
}
