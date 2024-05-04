package databases

import (
	"gorm.io/gorm"
	"talk_rater_bot/internal/data"
)

func AutoMigrateAllModels(db *gorm.DB) error {
	return db.AutoMigrate(&data.Conference{}, &data.User{}, &data.Lecture{}, &data.Evaluation{})
}
