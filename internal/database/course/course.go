package course

import (
	"errors"
	"fmt"
	"time"

	"github.com/niazlv/sport-plus-LCT/internal/config"
	"github.com/niazlv/sport-plus-LCT/internal/database/exercise"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Course модель курса
type Course struct {
	Id                int       `gorm:"primaryKey" json:"id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	Difficulty        string    `json:"difficulty"`
	DifficultyNumeric int       `json:"difficulty_numeric"`
	Direction         string    `json:"direction"`
	TrainerID         int       `json:"trainer_id"`
	Cost              float64   `json:"cost"`
	ParticipantsCount int       `json:"participants_count"`
	Rating            float64   `json:"rating"`
	RequiredTools     string    `json:"required_tools"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Classes           []Class   `json:"classes" gorm:"foreignKey:CourseID"`
}

// Class модель занятия
type Class struct {
	Id          int       `gorm:"primaryKey" json:"id"`
	CourseID    int       `json:"course_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Cover       string    `json:"cover"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Lessons     []Lesson  `json:"lessons" gorm:"foreignKey:ClassID"`
}

// Lesson модель урока
type Lesson struct {
	Id              int               `gorm:"primaryKey" json:"id"`
	CourseID        int               `json:"course_id"`
	ClassID         int               `json:"class_id"`
	ExerciseID      int               `json:"exercise_id"`
	Exercise        exercise.Exercise `json:"exercise" gorm:"foreignKey:ExerciseID"`
	DurationSeconds int               `json:"duration_seconds"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

// ClassImage модель изображения занятия
type ClassImage struct {
	Id       int    `gorm:"primaryKey" json:"id"`
	LessonID int    `json:"lesson_id"`
	Image    string `json:"image"`
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

	err = db.AutoMigrate(&Course{}, &Class{}, &Lesson{}, &ClassImage{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CRUD функции для модели Lesson

func CreateLesson(lesson *Lesson) (*Lesson, error) {
	result := db.Create(lesson)
	if result.Error != nil {
		return nil, result.Error
	}
	return lesson, nil
}

func GetLessonByID(id int) (*Lesson, error) {
	var lesson Lesson
	result := db.Preload("Images").Where("id = ?", id).First(&lesson)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &lesson, nil
}

func UpdateLesson(lesson *Lesson) error {
	result := db.Model(&Lesson{}).Where("id = ?", lesson.Id).Updates(lesson)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func DeleteLesson(id int) error {
	result := db.Delete(&Lesson{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// CRUD функции для модели Course

func CreateCourse(course *Course) (*Course, error) {
	result := db.Create(course)
	if result.Error != nil {
		return nil, result.Error
	}
	return course, nil
}

func GetCourseByID(id int) (*Course, error) {
	var course Course
	result := db.Preload("Classes").Where("id = ?", id).First(&course)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &course, nil
}

func UpdateCourse(course *Course) error {
	result := db.Model(&Course{}).Where("id = ?", course.Id).Updates(course)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func DeleteCourse(id int) error {
	result := db.Delete(&Course{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// CRUD функции для модели Class

func CreateClass(class *Class) (*Class, error) {
	result := db.Create(class)
	if result.Error != nil {
		return nil, result.Error
	}
	return class, nil
}

func GetClassByID(id int) (*Class, error) {
	var class Class
	result := db.Preload("Images").Where("id = ?", id).First(&class)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &class, nil
}

func UpdateClass(class *Class) error {
	result := db.Model(&Class{}).Where("id = ?", class.Id).Updates(class)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func DeleteClass(id int) error {
	result := db.Delete(&Class{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// CRUD функции для модели ClassImage

func CreateClassImage(classImage *ClassImage) (*ClassImage, error) {
	result := db.Create(classImage)
	if result.Error != nil {
		return nil, result.Error
	}
	return classImage, nil
}

func GetClassImageByID(id int) (*ClassImage, error) {
	var classImage ClassImage
	result := db.Where("id = ?", id).First(&classImage)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &classImage, nil
}

func UpdateClassImage(classImage *ClassImage) error {
	result := db.Model(&ClassImage{}).Where("id = ?", classImage.Id).Updates(classImage)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func DeleteClassImage(id int) error {
	result := db.Delete(&ClassImage{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
