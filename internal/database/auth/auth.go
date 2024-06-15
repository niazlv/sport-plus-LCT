package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/niazlv/sport-plus-LCT/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	TypeHeight = "height"
	TypeWeight = "weight"
	TypeWater  = "water"
)

var ValidTypesMeasurement = []string{TypeHeight, TypeWeight, TypeWater}

type Measurement struct {
	ID     int    `gorm:"primaryKey" json:"id"`
	UserID int    `json:"userId"`
	Date   string `json:"date" body:"date"`
	Value  string `json:"value" body:"value"`
	Type   string `json:"type" body:"type"`
}

type Train struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	UserID    int    `json:"userId"`
	Date      string `json:"date" body:"date"`
	TrainerID int    `json:"trainerId" body:"trainerId"`
	ClientID  int    `json:"clientId" body:"clientId"`
	Duration  string `json:"duration" body:"duration"`
	Trainer   User   `json:"trainer,omitempty" gorm:"foreignKey:TrainerID"`
	Client    User   `json:"client,omitempty" gorm:"foreignKey:ClientID"`
}

type User struct {
	Id               int           `gorm:"primaryKey" json:"id" body:"id"`
	Login            string        `gorm:"unique" json:"login" body:"login"`
	Password         string        `json:"password" body:"password"`
	Gender           string        `json:"gender" body:"gender"`
	Height           []Measurement `json:"height" body:"height" gorm:"foreignKey:UserID"`
	Weight           []Measurement `json:"weight" body:"weight" gorm:"foreignKey:UserID"`
	Water            []Measurement `json:"water" body:"water" gorm:"foreignKey:UserID"`
	Trains           []Train       `json:"trains" body:"trains" gorm:"foreignKey:UserID"`
	Goals            string        `json:"goals" body:"goals"`
	Experience       string        `json:"experience" body:"experience"`
	GymMember        bool          `json:"gymMember" body:"gymMember"`
	Beginner         bool          `json:"beginner" body:"beginner"`
	GymName          string        `json:"gymName" body:"gymName"`
	HealthConditions string        `json:"healthConditions" body:"healthConditions"`
	Role             int           `json:"role" body:"role"`
	Name             string        `json:"name" body:"name"`
	Icon             string        `json:"icon" body:"icon"`
	About            string        `json:"about"`
	Achivements      string        `json:"achivements"`
  Age              int           `json:"age"`
  Chats            []*Chat  `json:"chats" gorm:"many2many:chat_users"`
}

type Chat struct {
	Id        int       `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
  Users     []*User `gorm:"many2many:chat_users"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

	err = db.AutoMigrate(&User{}, &Measurement{}, &Train{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func FindUserByLogin(login string) (*User, error) {
	var user User
	result := db.Preload("Height", "type = ?", "height").
		Preload("Weight", "type = ?", "weight").
		Preload("Water", "type = ?", "water").
		Preload("Trains.Trainer").
		Preload("Trains.Client").
		Where("login = ?", login).First(&user)
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
	result := db.Preload("Height", "type = ?", "height").
		Preload("Weight", "type = ?", "weight").
		Preload("Water", "type = ?", "water").
		Preload("Trains.Trainer").
		Preload("Trains.Client").
		Where("id = ?", id).First(&user)
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

	if err := updateMeasurements(user); err != nil {
		return err
	}
	if err := updateTrains(user); err != nil {
		return err
	}

	return nil
}

func updateMeasurements(user *User) error {
	if err := updateHeight(user); err != nil {
		return err
	}
	if err := updateWeight(user); err != nil {
		return err
	}
	if err := updateWater(user); err != nil {
		return err
	}
	return nil
}

func updateHeight(user *User) error {
	db.Where("user_id = ? AND type = ?", user.Id, "height").Delete(&Measurement{})
	for _, measurement := range user.Height {
		measurement.UserID = user.Id
		db.Create(&measurement)
	}
	return nil
}

func updateWeight(user *User) error {
	db.Where("user_id = ? AND type = ?", user.Id, "weight").Delete(&Measurement{})
	for _, measurement := range user.Weight {
		measurement.UserID = user.Id
		db.Create(&measurement)
	}
	return nil
}

func updateWater(user *User) error {
	db.Where("user_id = ? AND type = ?", user.Id, "water").Delete(&Measurement{})
	for _, measurement := range user.Water {
		measurement.UserID = user.Id
		db.Create(&measurement)
	}
	return nil
}

func updateTrains(user *User) error {
	db.Where("user_id = ?", user.Id).Delete(&Train{})
	for _, train := range user.Trains {
		train.UserID = user.Id
		db.Create(&train)
	}
	return nil
}

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
	if user.Name != "" {
		updates["name"] = user.Name
	}
	if user.Icon != "" {
		updates["icon"] = user.Icon
	}
	if user.About != "" {
		updates["about"] = user.About
	}
	if user.Achivements != "" {
		updates["achivements"] = user.Achivements
	}
	if user.Age != 0 {
		updates["age"] = user.Age
	}

	result := db.Model(&User{}).Where("id = ?", user.Id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	if err := updateMeasurements(user); err != nil {
		return err
	}
	if err := updateTrains(user); err != nil {
		return err
	}

	return nil
}

func AddMeasurement(measurement *Measurement) (*Measurement, error) {
	result := db.Create(measurement)
	if result.Error != nil {
		return nil, result.Error
	}
	return measurement, nil
}

func UpdateMeasurement(measurement *Measurement) (*Measurement, error) {
	result := db.Save(measurement)
	if result.Error != nil {
		return nil, result.Error
	}
	return measurement, nil
}

func DeleteMeasurement(id int) error {
	result := db.Delete(&Measurement{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func AddTrain(train *Train) (*Train, error) {
	result := db.Create(train)
	if result.Error != nil {
		return nil, result.Error
	}
	// Подгружаем тренера и клиента после создания
	db.Preload("Trainer").Preload("Client").First(train, train.ID)
	return train, nil
}

func UpdateTrain(train *Train) (*Train, error) {
	result := db.Save(train)
	if result.Error != nil {
		return nil, result.Error
	}
	// Подгружаем тренера и клиента после обновления
	db.Preload("Trainer").Preload("Client").First(train, train.ID)
	return train, nil
}

func DeleteTrain(id int) error {
	result := db.Delete(&Train{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
