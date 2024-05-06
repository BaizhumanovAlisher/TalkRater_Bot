package admin

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type LectureInput struct {
	StartTime string // MSK time zone. Example of format: 21/07/2024 10:00:00
	Duration  string // in a file it is just a number. duration will be counted in minutes
	Title     string
	Speaker   string
	URL       string
}

func (ac *Controller) parseSchedule(pathFile string) ([]*LectureInput, error) {
	csvFile, err := os.Open(pathFile)
	if err != nil {
		return nil, err
	}

	// the number of reports described in tech requirements
	lectures := make([]*LectureInput, 0, 110)
	reader := csv.NewReader(csvFile)

	// program should ignore first line
	reader.Read()
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		if len(record) != 5 {
			return nil, fmt.Errorf("count of args should be 5, got %d.\nrecord: %s", len(record), record)
		}

		lecture := &LectureInput{
			StartTime: strings.TrimSpace(record[0]),
			Duration:  strings.TrimSpace(record[1]) + "m",
			Title:     strings.TrimSpace(record[2]),
			Speaker:   strings.TrimSpace(record[3]),
			URL:       strings.TrimSpace(record[4]),
		}

		lectures = append(lectures, lecture)
	}

	return lectures, nil
}
