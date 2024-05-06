package admin

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"talk_rater_bot/internal/data"
	"talk_rater_bot/internal/validators"
	"time"
)

func (ac *Controller) convertAndValidateLectures(lecturesInput []*LectureInput) ([]*data.Lecture, error) {
	lectures := make([]*data.Lecture, len(lecturesInput))
	// it is important get all errors from converting
	errorsSlice := make([]error, len(lectures))

	var wg sync.WaitGroup
	for i, lectureInput := range lecturesInput {
		wg.Add(1)
		go func(i int, lectureInput *LectureInput) {
			defer wg.Done()
			lecture, err := ac.convertAndValidateLecture(lectureInput)
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

	if count == 0 {
		return lectures, nil
	}

	output := result.String()
	return nil, errors.New(output[:len(output)-1])
}

func (ac *Controller) convertAndValidateLecture(input *LectureInput) (*data.Lecture, error) {
	startTime, err := ac.timeParser.ParseTimeString(input.StartTime)
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
	v.Check(ac.conference.StartTime.Before(lecture.Start) ||
		ac.conference.StartTime.Equal(lecture.Start),
		"conf start time", "lecture start time should be after conference start time or be equal")
	v.Check(ac.conference.EndTime.After(lecture.End) ||
		ac.conference.EndTime.Equal(lecture.End),
		"conf end time", "lecture end time should be after conference end time or be equal")

	if !v.Valid() {
		return nil, fmt.Errorf("error in line: %+v, description: %+v", input, v.Errors)
	}

	return lecture, nil
}
