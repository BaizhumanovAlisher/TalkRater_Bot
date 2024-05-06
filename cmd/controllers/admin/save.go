package admin

import (
	"gorm.io/gorm"
	"talk_rater_bot/internal/data"
	"time"
)

func (ac *Controller) save(lectures []*data.Lecture) error {
	return ac.db.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(lectures); i++ {
			result := ac.db.Save(lectures[i])

			if result.Error != nil {
				return result.Error
			}
		}

		return nil
	})
}

func (ac *Controller) extractSortedLectures() ([]*data.Lecture, error) {
	var lecturesInDB []*data.Lecture
	result := ac.db.Unscoped().Order("URL").Find(&lecturesInDB)

	if result.Error != nil {
		return nil, result.Error
	}

	return lecturesInDB, nil
}

func mergeLectures(oldLectures []*data.Lecture, newLectures []*data.Lecture) []*data.Lecture {
	outputLectures := make([]*data.Lecture, 0, len(oldLectures)+len(newLectures))
	i, j := 0, 0

	for i < len(newLectures) && j < len(oldLectures) {
		if newLectures[i].URL == oldLectures[j].URL {
			newLectures[i].ID = oldLectures[j].ID
			newLectures[i].DeletedAt.Valid = false

			outputLectures = append(outputLectures, newLectures[i])
			i++
			j++
		} else if newLectures[i].URL < oldLectures[j].URL {
			outputLectures = append(outputLectures, newLectures[i])
			i++
		} else {
			oldLectures[j].DeletedAt.Time = time.Now()
			oldLectures[j].DeletedAt.Valid = true
			outputLectures = append(outputLectures, oldLectures[j])

			j++
		}
	}

	for i < len(newLectures) {
		outputLectures = append(outputLectures, newLectures[i])
		i++
	}

	for j < len(oldLectures) {
		oldLectures[j].DeletedAt.Time = time.Now()
		oldLectures[j].DeletedAt.Valid = true
		outputLectures = append(outputLectures, oldLectures[j])

		j++
	}

	return outputLectures
}
