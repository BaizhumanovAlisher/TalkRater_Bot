package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"talk_rater_bot/internal/data"
	"talk_rater_bot/internal/validators"
	"time"
)

func (c *Controller) convertAndValidateLectures(lecturesInput []*LectureInput) ([]*data.Lecture, error) {
	err := checkUnique(lecturesInput)
	if err != nil {
		return nil, err
	}

	lectures := make([]*data.Lecture, len(lecturesInput))
	// it is important get all errors from converting
	errorsSlice := make([]error, len(lectures))

	var wg sync.WaitGroup
	for i, lectureInput := range lecturesInput {
		wg.Add(1)
		go func(i int, lectureInput *LectureInput) {
			defer wg.Done()
			lecture, err := c.convertAndValidateLecture(lectureInput)
			if err != nil {
				errorsSlice[i] = err
				return
			}

			lectures[i] = lecture
		}(i, lectureInput)
	}
	wg.Wait()

	count := 0
	var result bytes.Buffer
	for _, err := range errorsSlice {
		if err != nil {
			result.WriteString(err.Error())
			result.WriteString("\n")
			count++
		}
	}

	if count != 0 {
		output := result.String()
		return nil, errors.New(output[:len(output)-1])
	}

	return lectures, nil
}

func checkUnique(lecturesInput []*LectureInput) error {
	URLs := make([]string, len(lecturesInput))
	for i := 0; i < len(lecturesInput); i++ {
		URLs[i] = lecturesInput[i].URL
	}

	if !validators.Unique(URLs) {
		return errors.New("URLs must be unique")
	}

	return nil
}

func (c *Controller) convertAndValidateLecture(input *LectureInput) (*data.Lecture, error) {
	startTime, err := c.timeParser.ParseTimeString(input.StartTime)
	if err != nil {
		return nil, fmt.Errorf("error in line: %+v, description: %s", input, err.Error())
	}

	duration, err := time.ParseDuration(input.Duration)
	if err != nil {
		return nil, fmt.Errorf("error in line: %+v, description: %s", input, err.Error())
	}

	lecture := &data.Lecture{
		Title:   input.Title,
		Speaker: input.Speaker,
		URL:     input.URL,
		Start:   startTime,
		End:     startTime.Add(duration),
	}

	v := validators.New()
	data.ValidateLecture(v, lecture)
	v.Check(c.conference.StartTime.Before(lecture.Start) ||
		c.conference.StartTime.Equal(lecture.Start),
		"conf start time", "lecture start time should be after conference start time or be equal")
	v.Check(c.conference.EndTime.After(lecture.End) ||
		c.conference.EndTime.Equal(lecture.End),
		"conf end time", "lecture end time should be after conference end time or be equal")

	if !v.Valid() {
		return nil, fmt.Errorf("error in line: %+v, description: %+v", input, v.Errors)
	}

	return lecture, nil
}
