package exercise

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/niazlv/sport-plus-LCT/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Exercise модель занятия
type Exercise struct {
	Id               int       `gorm:"primaryKey" json:"id"`
	OriginalUri      string    `json:"original_uri"`
	Name             string    `json:"name"`
	Muscle           string    `json:"muscle"`
	AdditionalMuscle string    `json:"additional_muscle"`
	Type             string    `json:"type"`
	Equipment        string    `json:"equipment"`
	Difficulty       string    `json:"difficulty"`
	Duration         int       `json:"duration"`
	Photos           []Photo   `json:"photos" gorm:"foreignKey:ExerciseID"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Photo модель фотографии занятия
type Photo struct {
	Id         int    `gorm:"primaryKey" json:"id"`
	ExerciseID int    `json:"exercise_id"`
	URL        string `json:"url"`
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

	err = db.AutoMigrate(&Exercise{}, &Photo{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CRUD функции для модели Exercise

func CreateExercise(exercise *Exercise) (*Exercise, error) {
	result := db.Create(exercise)
	if result.Error != nil {
		return nil, result.Error
	}
	return exercise, nil
}

func GetExerciseByID(id int) (*Exercise, error) {
	var exercise Exercise
	result := db.Preload("Photos").Where("id = ?", id).First(&exercise)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	log.Println(exercise)
	return &exercise, nil
}

func UpdateExercise(exercise *Exercise) error {
	result := db.Model(&Exercise{}).Where("id = ?", exercise.Id).Updates(exercise)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func DeleteExercise(id int) error {
	result := db.Delete(&Exercise{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func FilterExercises(additionalMuscle string) ([]Exercise, error) {
	var exercises []Exercise
	result := db.Preload("Photos").Where("additional_muscle LIKE ?", "%"+additionalMuscle+"%").Find(&exercises)
	if result.Error != nil {
		return nil, result.Error
	}
	return exercises, nil
}
