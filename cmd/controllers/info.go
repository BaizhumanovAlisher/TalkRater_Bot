package controllers

import "talk_rater_bot/internal/data"

func (c *Controller) GetCurrentConference() *data.Conference {
	tmp := *c.conference
	return &tmp
}

func (c *Controller) GetSchedule(limit, offset int) ([]*data.Lecture, error) {
	var lectures []*data.Lecture
	result := c.db.Order("start").Limit(limit).Offset(offset).Find(&lectures)

	if result.Error != nil {
		return nil, result.Error
	}

	return lectures, nil
}

func (c *Controller) CountPageInLectures(pageSize int64) (int64, error) {
	var count int64
	result := c.db.Model(&data.Lecture{}).Count(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return (count + pageSize - 1) / pageSize, nil
}

func (c *Controller) GetLecture(id int64) (*data.Lecture, error) {
	var lecture data.Lecture
	result := c.db.First(&lecture, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &lecture, nil
}
