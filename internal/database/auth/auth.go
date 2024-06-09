package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/niazlv/sport-plus-LCT/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	Id               int    `gorm:"primaryKey"`
	Login            string `gorm:"unique" json:"login"`
	Password         string `json:"password"`
	Gender           string `json:"gender"`
	Height           int    `json:"height"`
	Weight           int    `json:"weight"`
	Goals            string `json:"goals"`
	Experience       string `json:"experience"`
	GymMember        bool   `json:"gymMember"`
	Beginner         bool   `json:"beginner"`
	GymName          string `json:"gymName"`
	HealthConditions string `json:"healthConditions"`
	Role             int    `json:"role"`
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
			// Если произошла ошибка при подключении к базе данных, ждем некоторое время перед повторной попыткой
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	if db == nil {
		return nil, errors.New("failed to connect to database")
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func FindUserByLogin(login string) (*User, error) {
	var user User
	result := db.Where("login = ?", login).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func CreateUser(user *User) (*User, error) {
	result := db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func FindUserByID(id int) (*User, error) {
	var user User
	result := db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func UpdateUserName(userId int, newName string) error {
	result := db.Model(&User{}).Where("id = ?", userId).Update("name", newName)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
