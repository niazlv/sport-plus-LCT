package calendar

import (
	"errors"
	"fmt"
	"time"

	"github.com/niazlv/sport-plus-LCT/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Schedule struct {
	Id             int       `gorm:"primaryKey" json:"id"`
	CoachID        int       `json:"coach_id"`  // ID тренера
	ClientID       int       `json:"client_id"` // ID клиента
	Date           time.Time `json:"date"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	Type           string    `json:"type"`
	ReminderClient bool      `json:"reminder_client"`
	ReminderCoach  bool      `json:"reminder_coach"`
	IsGlobal       bool      `json:"is_global"` // Глобальное или локальное событие
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

var db *gorm.DB

func InitDB() (*gorm.DB, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort)

	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	if db == nil {
		return nil, errors.New("failed to connect to database")
	}

	err = db.AutoMigrate(&Schedule{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateSchedule(schedule *Schedule) (*Schedule, error) {
	result := db.Create(schedule)
	if result.Error != nil {
		return nil, result.Error
	}
	return schedule, nil
}

func GetScheduleByID(id int) (*Schedule, error) {
	var schedule Schedule
	result := db.Where("id = ?", id).First(&schedule)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &schedule, nil
}

func UpdateSchedule(schedule *Schedule) error {
	result := db.Model(&Schedule{}).Where("id = ?", schedule.Id).Updates(schedule)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func DeleteSchedule(id int) error {
	result := db.Delete(&Schedule{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func GetSchedulesByCoachID(coachID int) ([]Schedule, error) {
	var schedules []Schedule
	result := db.Where("coach_id = ? OR is_global = ?", coachID, true).Find(&schedules)
	if result.Error != nil {
		return nil, result.Error
	}
	return schedules, nil
}
