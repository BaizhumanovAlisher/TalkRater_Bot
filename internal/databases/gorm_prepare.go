package databases

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"os"
	"os/exec"
	"strconv"
	"talk_rater_bot/internal/data"
	"talk_rater_bot/internal/helpers"
	"time"
)

type PrepareDBHelper struct {
	db                        *gorm.DB
	currConference            *data.Conference
	cleanupDBForNewConference bool
	backupOptions             *helpers.DatabaseConfig
	pathToBackup              string
}

func NewPrepareDBHelper(db *gorm.DB, currConference *data.Conference, clearDbForNewConference bool, backupOptions *helpers.DatabaseConfig, pathToBackup string) *PrepareDBHelper {
	return &PrepareDBHelper{db: db, currConference: currConference, cleanupDBForNewConference: clearDbForNewConference, backupOptions: backupOptions, pathToBackup: pathToBackup}
}

func (ph *PrepareDBHelper) PrepareDB() error {
	err := ph.db.AutoMigrate(&data.Conference{}, &data.User{}, &data.Lecture{}, &data.Evaluation{})
	if err != nil {
		return err
	}

	if ph.cleanupDBForNewConference {
		return ph.updateConference()
	}

	return nil
}

func (ph *PrepareDBHelper) updateConference() error {
	var conferences []data.Conference
	result := ph.db.Find(&conferences)

	if result.Error != nil {
		return result.Error
	}

	if len(conferences) > 1 {
		return errors.New("count of conferences can not be greater than 1. Incorrect db data")
	} else if len(conferences) == 1 {
		prevConference := &conferences[0]
		if data.AreEqualConferences(prevConference, ph.currConference) {
			return nil
		}

		err := ph.backupDatabase()
		if err != nil {
			return err
		}

		err = ph.cleanupDatabase()
		if err != nil {
			return err
		}
	}

	return ph.db.Create(ph.currConference).Error
}

func (ph *PrepareDBHelper) backupDatabase() error {
	backupFile := fmt.Sprintf("%s%sbackup_%s.sql", ph.pathToBackup, string(os.PathSeparator), time.Now().Format("2006-01-02_15-04-05"))

	cmd := exec.Command("pg_dump",
		"-h", ph.backupOptions.Host,
		"-p", strconv.Itoa(ph.backupOptions.Port),
		"-U", ph.backupOptions.User,
		"-d", ph.backupOptions.DatabaseName,
		"-f", backupFile,
	)

	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", ph.backupOptions.Password))

	err := cmd.Run()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return fmt.Errorf("pg_dump exited with status %d: %v", exitError.ExitCode(), err)
		}

		return fmt.Errorf("failed to run pg_dump: %v", err)
	}

	return nil
}

func (ph *PrepareDBHelper) cleanupDatabase() error {
	result := ph.db.Unscoped().Where("1 = 1").Delete(&data.Evaluation{})
	if result.Error != nil {
		return result.Error
	}

	result = ph.db.Unscoped().Where("1 = 1").Delete(&data.User{})
	if result.Error != nil {
		return result.Error
	}

	result = ph.db.Unscoped().Where("1 = 1").Delete(&data.Lecture{})
	if result.Error != nil {
		return result.Error
	}

	result = ph.db.Unscoped().Where("1 = 1").Delete(&data.Conference{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
