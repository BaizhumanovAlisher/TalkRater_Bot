package controllers

import "talk_rater_bot/internal/data"

func (c *Controller) GetCurrentConference() *data.Conference {
	tmp := *c.conference
	return &tmp
}

func (c *Controller) GetSchedule() ([]*data.Lecture, error) {
	var lecturers []*data.Lecture
	result := c.db.Order("start").Find(&lecturers)

	if result.Error != nil {
		return nil, result.Error
	}

	aLotOf := make([]*data.Lecture, 0, len(lecturers)*25)

	for i := 0; i < 5; i++ {
		for _, lect := range lecturers {
			aLotOf = append(aLotOf, lect)
		}
	}

	return aLotOf, nil
}
