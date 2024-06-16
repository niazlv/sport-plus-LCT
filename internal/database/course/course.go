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

const (
	StatusNotStarted = "Не начато"
	StatusInProgress = "В процессе"
	StatusCompleted  = "Завершено"
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
	Id              int              `gorm:"primaryKey" json:"id"`
	CourseID        int              `json:"course_id"`
	ClassID         int              `json:"class_id"`
	DurationSeconds int              `json:"duration_seconds"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	Exercises       []LessonExercise `json:"exercises" gorm:"foreignKey:LessonID"`
}

// LessonExercise модель связи между уроком и упражнением
type LessonExercise struct {
	Id         int               `gorm:"primaryKey" json:"id"`
	LessonID   int               `json:"lesson_id"`
	ExerciseID int               `json:"exercise_id"`
	Exercise   exercise.Exercise `json:"exercise" gorm:"foreignKey:ExerciseID"`
}

// ClassImage модель изображения занятия
type ClassImage struct {
	Id       int    `gorm:"primaryKey" json:"id"`
	LessonID int    `json:"lesson_id"`
	Image    string `json:"image"`
}

// Progress status
type ExerciseStatus struct {
	Id         int    `gorm:"primaryKey" json:"id"`
	LessonID   int    `json:"lesson_id"`   // Foreign key to LessonStatus
	ExerciseID int    `json:"exercise_id"` // Foreign key to Exercise
	Status     string `json:"status"`
}

type LessonStatus struct {
	Id        int              `gorm:"primaryKey" json:"id"`
	ClassID   int              `json:"class_id"`  // Foreign key to ClassStatus
	LessonID  int              `json:"lesson_id"` // Foreign key to Lesson
	Status    string           `json:"status"`
	Exercises []ExerciseStatus `json:"exercises" gorm:"foreignKey:LessonID"`
}

type ClassStatus struct {
	Id       int            `gorm:"primaryKey" json:"id"`
	CourseID int            `json:"course_id"` // Foreign key to CourseStatus
	ClassID  int            `json:"class_id"`  // Foreign key to Class
	Status   string         `json:"status"`
	Lessons  []LessonStatus `json:"lessons" gorm:"foreignKey:ClassID"`
}

type CourseStatus struct {
	Id       int           `gorm:"primaryKey" json:"id"`
	ClientID int           `json:"client_id"` // Foreign key to Client
	CourseID int           `json:"course_id"` // Foreign key to Course
	Status   string        `json:"status"`
	Classes  []ClassStatus `json:"classes" gorm:"foreignKey:CourseID"`
}

type ClientProgress struct {
	Id        int            `gorm:"primaryKey" json:"id"`
	ClientID  int            `json:"client_id"` // Foreign key to Client
	Courses   []CourseStatus `json:"courses" gorm:"foreignKey:ClientID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
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

	err = db.AutoMigrate(&Course{}, &Class{}, &Lesson{}, &ClassImage{}, &LessonExercise{}, &ClientProgress{}, &CourseStatus{}, &ClassStatus{}, &LessonStatus{}, &ExerciseStatus{})
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
	result := db.Preload("Exercises.Exercise.Photos").Where("id = ?", id).First(&lesson)
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

	// Обновление упражнений
	db.Where("lesson_id = ?", lesson.Id).Delete(&LessonExercise{})
	for _, exercise := range lesson.Exercises {
		exercise.LessonID = lesson.Id
		db.Create(&exercise)
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
	db.Where("lesson_id = ?", id).Delete(&LessonExercise{})
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
	result := db.Preload("Classes").Preload("Classes.Lessons").Where("id = ?", id).First(&course)
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

// CRUD функции для модели ClientProgress

func CreateClientProgress(progress *ClientProgress) (*ClientProgress, error) {
	result := db.Create(progress)
	if result.Error != nil {
		return nil, result.Error
	}
	return progress, nil
}

func GetClientProgressByID(id int) (*ClientProgress, error) {
	var progress ClientProgress
	result := db.Preload("Courses.Classes.Lessons.Exercises").Where("id = ?", id).First(&progress)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &progress, nil
}

func GetClientProgressByClientAndCourseID(clientID int, courseID int) (*ClientProgress, error) {
	var progress ClientProgress
	result := db.Preload("Courses.Classes.Lessons.Exercises").
		Where("client_id = ?", clientID).
		First(&progress)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}

	// Ensure the structure is fully populated
	EnsureFullStructure(clientID, courseID, &progress)

	return &progress, nil
}

func EnsureFullStructure(clientID int, courseID int, progress *ClientProgress) {
	var course Course
	db.Preload("Classes.Lessons.Exercises").First(&course, courseID)

	for _, class := range course.Classes {
		classFound := false
		for _, classStatus := range progress.Courses[0].Classes {
			if classStatus.ClassID == class.Id {
				classFound = true
				break
			}
		}
		if !classFound {
			newClassStatus := ClassStatus{
				CourseID: courseID,
				ClassID:  class.Id,
				Status:   StatusNotStarted,
				Lessons:  []LessonStatus{},
			}
			progress.Courses[0].Classes = append(progress.Courses[0].Classes, newClassStatus)
		}

		for _, lesson := range class.Lessons {
			lessonFound := false
			for _, lessonStatus := range progress.Courses[0].Classes[0].Lessons {
				if lessonStatus.LessonID == lesson.Id {
					lessonFound = true
					break
				}
			}
			if !lessonFound {
				newLessonStatus := LessonStatus{
					ClassID:   class.Id,
					LessonID:  lesson.Id,
					Status:    StatusNotStarted,
					Exercises: []ExerciseStatus{},
				}
				progress.Courses[0].Classes[0].Lessons = append(progress.Courses[0].Classes[0].Lessons, newLessonStatus)
			}

			for _, exercise := range lesson.Exercises {
				exerciseFound := false
				for _, exerciseStatus := range progress.Courses[0].Classes[0].Lessons[0].Exercises {
					if exerciseStatus.ExerciseID == exercise.Id {
						exerciseFound = true
						break
					}
				}
				if !exerciseFound {
					newExerciseStatus := ExerciseStatus{
						LessonID:   lesson.Id,
						ExerciseID: exercise.Id,
						Status:     StatusNotStarted,
					}
					progress.Courses[0].Classes[0].Lessons[0].Exercises = append(progress.Courses[0].Classes[0].Lessons[0].Exercises, newExerciseStatus)
				}
			}
		}
	}
}

func EnsureFullClientStructure(clientID int, progress *ClientProgress) {
	var courses []Course
	db.Preload("Classes.Lessons.Exercises").Find(&courses)

	for _, course := range courses {
		var courseStatus *CourseStatus
		courseFound := false
		for i := range progress.Courses {
			if progress.Courses[i].CourseID == course.Id {
				courseStatus = &progress.Courses[i]
				courseFound = true
				break
			}
		}
		if !courseFound {
			newCourseStatus := CourseStatus{
				ClientID: clientID,
				CourseID: course.Id,
				Status:   StatusNotStarted,
				Classes:  []ClassStatus{},
			}
			// Сохраняем новый CourseStatus в базу данных
			db.Create(&newCourseStatus)
			progress.Courses = append(progress.Courses, newCourseStatus)
			courseStatus = &progress.Courses[len(progress.Courses)-1]
		}

		for _, class := range course.Classes {
			var classStatus *ClassStatus
			classFound := false
			for i := range courseStatus.Classes {
				if courseStatus.Classes[i].ClassID == class.Id {
					classStatus = &courseStatus.Classes[i]
					classFound = true
					break
				}
			}
			if !classFound {
				newClassStatus := ClassStatus{
					CourseID: courseStatus.Id, // Используем реальный CourseID
					ClassID:  class.Id,
					Status:   StatusNotStarted,
					Lessons:  []LessonStatus{},
				}
				// Сохраняем новый ClassStatus в базу данных
				db.Create(&newClassStatus)
				courseStatus.Classes = append(courseStatus.Classes, newClassStatus)
				classStatus = &courseStatus.Classes[len(courseStatus.Classes)-1]
			}

			for _, lesson := range class.Lessons {
				var lessonStatus *LessonStatus
				lessonFound := false
				for i := range classStatus.Lessons {
					if classStatus.Lessons[i].LessonID == lesson.Id {
						lessonStatus = &classStatus.Lessons[i]
						lessonFound = true
						break
					}
				}
				if !lessonFound {
					newLessonStatus := LessonStatus{
						ClassID:   classStatus.Id, // Используем реальный ClassID
						LessonID:  lesson.Id,
						Status:    StatusNotStarted,
						Exercises: []ExerciseStatus{},
					}
					// Сохраняем новый LessonStatus в базу данных
					db.Create(&newLessonStatus)
					classStatus.Lessons = append(classStatus.Lessons, newLessonStatus)
					lessonStatus = &classStatus.Lessons[len(classStatus.Lessons)-1]
				}

				for _, exercise := range lesson.Exercises {
					exerciseFound := false
					for _, exerciseStatus := range lessonStatus.Exercises {
						if exerciseStatus.ExerciseID == exercise.Id {
							exerciseFound = true
							break
						}
					}
					if !exerciseFound {
						newExerciseStatus := ExerciseStatus{
							LessonID:   lessonStatus.Id, // Используем реальный LessonID
							ExerciseID: exercise.Id,
							Status:     StatusNotStarted,
						}
						// Сохраняем новый ExerciseStatus в базу данных
						db.Create(&newExerciseStatus)
						lessonStatus.Exercises = append(lessonStatus.Exercises, newExerciseStatus)
					}
				}
			}
		}
	}
}

func UpdateClientProgress(progress *ClientProgress) error {
	result := db.Save(progress)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteClientProgress(id int) error {
	result := db.Delete(&ClientProgress{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func UpdateCourseStatus(clientID int, courseID int, newStatus string) error {
	var courseStatus CourseStatus
	result := db.Where("client_id = ? AND course_id = ?", clientID, courseID).First(&courseStatus)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			newCourseStatus := CourseStatus{
				ClientID: clientID,
				CourseID: courseID,
				Status:   newStatus,
				Classes:  []ClassStatus{},
			}
			return db.Create(&newCourseStatus).Error
		}
		return result.Error
	}
	courseStatus.Status = newStatus
	return db.Save(&courseStatus).Error
}

func UpdateClassStatus(clientID int, courseID int, classID int, newStatus string) error {
	var classStatus ClassStatus
	result := db.Where("course_id = ? AND class_id = ?", courseID, classID).First(&classStatus)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			newClassStatus := ClassStatus{
				CourseID: courseID,
				ClassID:  classID,
				Status:   newStatus,
				Lessons:  []LessonStatus{},
			}
			return db.Create(&newClassStatus).Error
		}
		return result.Error
	}
	classStatus.Status = newStatus
	return db.Save(&classStatus).Error
}

func UpdateLessonStatus(clientID int, classID int, lessonID int, newStatus string) error {
	var lessonStatus LessonStatus
	result := db.Where("class_id = ? AND lesson_id = ?", classID, lessonID).First(&lessonStatus)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			newLessonStatus := LessonStatus{
				ClassID:   classID,
				LessonID:  lessonID,
				Status:    newStatus,
				Exercises: []ExerciseStatus{},
			}
			return db.Create(&newLessonStatus).Error
		}
		return result.Error
	}
	lessonStatus.Status = newStatus
	return db.Save(&lessonStatus).Error
}

func UpdateExerciseStatus(clientID int, lessonID int, exerciseID int, newStatus string) error {
	var exerciseStatus ExerciseStatus
	result := db.Where("lesson_id = ? AND exercise_id = ?", lessonID, exerciseID).First(&exerciseStatus)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			newExerciseStatus := ExerciseStatus{
				LessonID:   lessonID,
				ExerciseID: exerciseID,
				Status:     newStatus,
			}
			return db.Create(&newExerciseStatus).Error
		}
		return result.Error
	}
	exerciseStatus.Status = newStatus
	return db.Save(&exerciseStatus).Error
}
func GetClientProgressByClientID(clientID int) (*ClientProgress, error) {
	var progress ClientProgress
	result := db.Preload("Courses.Classes.Lessons.Exercises").Where("client_id = ?", clientID).First(&progress)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &progress, nil
}
