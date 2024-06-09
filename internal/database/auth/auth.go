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
	Id               int    `gorm:"primaryKey" json:"id" body:"id"`
	Login            string `gorm:"unique" json:"login" body:"login"`
	Password         string `json:"password" body:"password"`
	Gender           string `json:"gender" body:"gender"`
	Height           int    `json:"height" body:"height"`
	Weight           int    `json:"weight" body:"weight"`
	Goals            string `json:"goals" body:"goals"`
	Experience       string `json:"experience" body:"experience"`
	GymMember        bool   `json:"gymMember" body:"gymMember"`
	Beginner         bool   `json:"beginner" body:"beginner"`
	GymName          string `json:"gymName" body:"gymName"`
	HealthConditions string `json:"healthConditions" body:"healthConditions"`
	Role             int    `json:"role" body:"role"`
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

func UpdateUser(user *User) error {
	result := db.Model(&User{}).Where("id = ?", user.Id).Updates(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// PartialUpdateUser updates the user record with the provided data, ignoring zero-value fields
func PartialUpdateUser(user *User) error {
	updates := make(map[string]interface{})

	if user.Login != "" {
		updates["login"] = user.Login
	}
	if user.Password != "" {
		updates["password"] = user.Password
	}
	if user.Gender != "" {
		updates["gender"] = user.Gender
	}
	if user.Height != 0 {
		updates["height"] = user.Height
	}
	if user.Weight != 0 {
		updates["weight"] = user.Weight
	}
	if user.Goals != "" {
		updates["goals"] = user.Goals
	}
	if user.Experience != "" {
		updates["experience"] = user.Experience
	}
	if user.GymMember {
		updates["gym_member"] = user.GymMember
	}
	if user.Beginner {
		updates["beginner"] = user.Beginner
	}
	if user.GymName != "" {
		updates["gym_name"] = user.GymName
	}
	if user.HealthConditions != "" {
		updates["health_conditions"] = user.HealthConditions
	}
	if user.Role != 0 {
		updates["role"] = user.Role
	}

	result := db.Model(&User{}).Where("id = ?", user.Id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
