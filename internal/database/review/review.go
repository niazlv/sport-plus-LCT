package review

import (
	"errors"
	"fmt"
	"time"

	"github.com/niazlv/sport-plus-LCT/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Review модель отзыва
type Review struct {
	Id               int       `gorm:"primaryKey" json:"id"`
	ClassID          int       `json:"class_id"`
	ClientID         int       `json:"client_id"`
	TrainerID        int       `json:"trainer_id"`
	DifficultyRating int       `json:"difficulty_rating"`
	WellBeingRating  int       `json:"well_being_rating"`
	OverallRating    int       `json:"overall_rating"`
	Comment          string    `json:"comment"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
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

	err = db.AutoMigrate(&Review{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CRUD функции для модели Review

func CreateReview(review *Review) (*Review, error) {
	result := db.Create(review)
	if result.Error != nil {
		return nil, result.Error
	}
	return review, nil
}

func GetReviewByID(id int) (*Review, error) {
	var review Review
	result := db.First(&review, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &review, nil
}

func GetReviewsByClassID(classID int) ([]Review, error) {
	var reviews []Review
	result := db.Where("class_id = ?", classID).Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}
	return reviews, nil
}

func UpdateReview(review *Review) error {
	result := db.Save(review)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteReview(id int) error {
	result := db.Delete(&Review{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
