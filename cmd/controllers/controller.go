package controllers

import (
	"gorm.io/gorm"
	"sort"
	"sync"
	"talk_rater_bot/internal/data"
	"talk_rater_bot/internal/helpers"
)

type Controller struct {
	db         *gorm.DB
	timeParser *helpers.TimeParser
	conference *data.Conference
}

func NewController(DB *gorm.DB, timeParser *helpers.TimeParser, conference *data.Conference) *Controller {
	return &Controller{db: DB, timeParser: timeParser, conference: conference}
}

func (c *Controller) GenerateSchedule(pathFile string) error {
	lecturesInput, err := c.parseSchedule(pathFile)
	if err != nil {
		return err
	}

	newLectures, err := c.convertAndValidateLectures(lecturesInput)
	if err != nil {
		return err
	}

	oldLectures, err := sortAndReadFromDBConcurrently(newLectures, c)

	mergedLectures := mergeLectures(oldLectures, newLectures)

	return c.save(mergedLectures)
}

func sortAndReadFromDBConcurrently(newLectures []*data.Lecture, ac *Controller) (oldLectures []*data.Lecture, err error) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		sort.Slice(newLectures, func(i, j int) bool {
			return newLectures[i].URL < newLectures[j].URL
		})
	}()

	go func() {
		defer wg.Done()
		oldLectures, err = ac.extractSortedLectures()
	}()
	wg.Wait()

	if err != nil {
		return nil, err
	}

	return oldLectures, nil
}

func (c *Controller) ExportEvaluations() ([]*data.ExportEvaluation, error) {
	var exportEvaluations []*data.ExportEvaluation

	result := c.db.Table("evaluations").
		Select("users.identity_info as user, lectures.url as url, evaluations.score_content as content, evaluations.score_performance as performance, evaluations.comment as comment").
		Joins("left join users on users.id = evaluations.user_id").
		Joins("left join lectures on lectures.id = evaluations.lecture_id").
		Where("evaluations.type_evaluation = ?", string(data.Correct)).
		Scan(&exportEvaluations)

	if result.Error != nil {
		return nil, result.Error
	}

	if exportEvaluations == nil {
		exportEvaluations = make([]*data.ExportEvaluation, 0)
	}

	return exportEvaluations, nil
}
