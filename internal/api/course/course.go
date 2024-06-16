package course

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/niazlv/sport-plus-LCT/internal/api/auth"
	"github.com/niazlv/sport-plus-LCT/internal/database/course"
	database_course "github.com/niazlv/sport-plus-LCT/internal/database/course"
	"github.com/wI2L/fizz"
	"gorm.io/gorm"
)

var db *gorm.DB

func Setup(rg *fizz.RouterGroup) {
	api := rg.Group("course", "Course", "Course related endpoints")

	var err error
	db, err = course.InitDB()
	if err != nil {
		log.Fatal("db courses can't be init: ", err)
	}

	api.GET("", []fizz.OperationOption{fizz.Summary("Get list of courses"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetCourses, 200))
	api.GET("/:course_id", []fizz.OperationOption{fizz.Summary("Get course by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetCourseByID, 200))
	api.POST("", []fizz.OperationOption{fizz.Summary("Create a new course"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(CreateCourse, 201))
	api.PUT("/:course_id", []fizz.OperationOption{fizz.Summary("Update course by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(UpdateCourse, 200))

	api.GET("/progress", []fizz.OperationOption{fizz.Summary("Get full client progress"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetFullClientProgress, 200))
	api.GET("/progress/:course_id", []fizz.OperationOption{fizz.Summary("Get course progress by course ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetCourseProgress, 200))
	api.PUT("/progress/:course_id", []fizz.OperationOption{fizz.Summary("Update course progress by course ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(UpdateCourseProgress, 200))
	api.PUT("/progress/:course_id/class/:class_id", []fizz.OperationOption{fizz.Summary("Update class progress by class ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(UpdateClassProgress, 200))
	api.PUT("/progress/:course_id/class/:class_id/lesson/:lesson_id", []fizz.OperationOption{fizz.Summary("Update lesson progress by lesson ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(UpdateLessonProgress, 200))
	api.PUT("/progress/:course_id/class/:class_id/lesson/:lesson_id/exercise/:exercise_id", []fizz.OperationOption{fizz.Summary("Update exercise progress by exercise ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(UpdateExerciseProgress, 200))

	SetupClassRoutes(api)
}

type CourseOutput struct {
	Course course.Course `json:"course"`
}

type CoursesOutput struct {
	Courses []course.Course `json:"courses"`
}

type GetCourseByIDParams struct {
	ID string `path:"course_id" binding:"required"`
}

func GetCourses(c *gin.Context) (*CoursesOutput, error) {
	log.Println("GetCourses called")

	var courses []course.Course
	result := db.Find(&courses)
	if result.Error != nil {
		log.Println("Error retrieving courses:", result.Error)
		return nil, result.Error
	}

	log.Printf("Retrieved courses: %+v\n", courses)
	return &CoursesOutput{
		Courses: courses,
	}, nil
}

func GetCourseByID(c *gin.Context, params *GetCourseByIDParams) (*CourseOutput, error) {
	idStr := params.ID
	log.Println("GetCourseByID called with ID:", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid course ID:", idStr)
		return nil, &gin.Error{
			Err:  errors.New("invalid course_id"),
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "invalid course_id"},
		}
	}

	var course *course.Course
	course, err = database_course.GetCourseByID(id)
	if err != nil {
		log.Println("Error retrieving course:", err)
		if err == gorm.ErrRecordNotFound {
			return nil, &gin.Error{
				Err:  err,
				Type: gin.ErrorTypePublic,
				Meta: gin.H{"error": "course not found"},
			}
		}
		return nil, err
	}

	log.Printf("Retrieved course: %+v\n", course)
	return &CourseOutput{
		Course: *course,
	}, nil
}

type CreateCourseInput struct {
	Title             string  `json:"title" binding:"required"`
	Description       string  `json:"description"`
	Difficulty        string  `json:"difficulty"`
	DifficultyNumeric int     `json:"difficulty_numeric"`
	Direction         string  `json:"direction"`
	TrainerID         int     `json:"trainer_id"`
	Cost              float64 `json:"cost"`
	ParticipantsCount int     `json:"participants_count"`
	Rating            float64 `json:"rating"`
	RequiredTools     string  `json:"required_tools"`
}

func CreateCourse(c *gin.Context, in *CreateCourseInput) (*CourseOutput, error) {
	log.Printf("CreateCourse called with input: %+v\n", in)

	newCourse := course.Course{
		Title:             in.Title,
		Description:       in.Description,
		Difficulty:        in.Difficulty,
		DifficultyNumeric: in.DifficultyNumeric,
		Direction:         in.Direction,
		TrainerID:         in.TrainerID,
		Cost:              in.Cost,
		ParticipantsCount: in.ParticipantsCount,
		Rating:            in.Rating,
		RequiredTools:     in.RequiredTools,
	}

	result := db.Create(&newCourse)
	if result.Error != nil {
		log.Println("Error creating course:", result.Error)
		return nil, result.Error
	}

	log.Printf("Created course: %+v\n", newCourse)
	return &CourseOutput{
		Course: newCourse,
	}, nil
}

type UpdateCourseInput struct {
	ID                string  `path:"course_id"`
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	Difficulty        string  `json:"difficulty"`
	DifficultyNumeric int     `json:"difficulty_numeric"`
	Direction         string  `json:"direction"`
	TrainerID         int     `json:"trainer_id"`
	Cost              float64 `json:"cost"`
	ParticipantsCount int     `json:"participants_count"`
	Rating            float64 `json:"rating"`
	RequiredTools     string  `json:"required_tools"`
}

func UpdateCourse(c *gin.Context, in *UpdateCourseInput) (*CourseOutput, error) {
	id := in.ID
	log.Printf("UpdateCourse called with ID: %s and input: %+v\n", id, in)

	var course course.Course
	// result := db.First(&course, id)

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	courses, err := database_course.GetCourseByID(idInt)
	course = *courses
	if err != nil {
		log.Println("Error retrieving course:", err)
		if err == gorm.ErrRecordNotFound {
			return nil, &gin.Error{
				Err:  err,
				Type: gin.ErrorTypePublic,
				Meta: gin.H{"error": "course not found"},
			}
		}
		return nil, err
	}

	if in.Title != "" {
		course.Title = in.Title
	}
	if in.Description != "" {
		course.Description = in.Description
	}
	if in.Difficulty != "" {
		course.Difficulty = in.Difficulty
	}
	if in.DifficultyNumeric != 0 {
		course.DifficultyNumeric = in.DifficultyNumeric
	}
	if in.Direction != "" {
		course.Direction = in.Direction
	}
	if in.TrainerID != 0 {
		course.TrainerID = in.TrainerID
	}
	if in.Cost != 0 {
		course.Cost = in.Cost
	}
	if in.ParticipantsCount != 0 {
		course.ParticipantsCount = in.ParticipantsCount
	}
	if in.Rating != 0 {
		course.Rating = in.Rating
	}
	if in.RequiredTools != "" {
		course.RequiredTools = in.RequiredTools
	}

	result := db.Save(&course)
	if result.Error != nil {
		log.Println("Error updating course:", result.Error)
		return nil, result.Error
	}

	log.Printf("Updated course: %+v\n", course)
	return &CourseOutput{
		Course: course,
	}, nil
}

// Endpoint для получения прогресса курса
type GetCourseProgressParams struct {
	CourseID string `path:"course_id" binding:"required"`
}

func GetCourseProgress(c *gin.Context, params *GetCourseProgressParams) (*course.CourseStatus, error) {
	courseID, err := strconv.Atoi(params.CourseID)
	if err != nil {
		return nil, &gin.Error{
			Err:  errors.New("invalid course_id"),
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "invalid course_id"},
		}
	}

	claims := c.MustGet("claims").(jwt.MapClaims)
	userClaims, err := auth.ExtractClaims(claims)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	progress, err := course.GetClientProgressByClientAndCourseID(userClaims.ID, courseID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if progress == nil {
		// Создаем новую запись прогресса, если не существует
		newProgress := course.ClientProgress{
			ClientID: userClaims.ID,
			Courses: []course.CourseStatus{
				{
					ClientID: userClaims.ID,
					CourseID: courseID,
					Status:   course.StatusNotStarted,
					Classes:  []course.ClassStatus{},
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		progress, err = course.CreateClientProgress(&newProgress)
		if err != nil {
			return nil, err
		}
	}

	// Ensure the structure is fully populated for this specific course
	database_course.EnsureFullStructure(userClaims.ID, courseID, progress)

	for _, courseStatus := range progress.Courses {
		if courseStatus.CourseID == courseID {
			return &courseStatus, nil
		}
	}

	return nil, gorm.ErrRecordNotFound
}

// Endpoint для обновления прогресса курса
type UpdateCourseProgressInput struct {
	Status   string `json:"status" binding:"required"`
	CourseID string `path:"course_id"`
}

func UpdateCourseProgress(c *gin.Context, in *UpdateCourseProgressInput) (*course.ClientProgress, error) {
	courseID, err := strconv.Atoi(in.CourseID)
	if err != nil {
		return nil, &gin.Error{
			Err:  errors.New("invalid course_id"),
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "invalid course_id"},
		}
	}

	claims := c.MustGet("claims").(jwt.MapClaims)

	userClaims, err := auth.ExtractClaims(claims)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	err = course.UpdateCourseStatus(userClaims.ID, courseID, in.Status)
	if err != nil {
		return nil, err
	}

	progress, err := course.GetClientProgressByClientAndCourseID(userClaims.ID, courseID)
	if err != nil {
		return nil, err
	}

	return progress, nil
}

// type UpdateProgressParams struct {
// 	CourseID string `path:"course_id" binding:"required"`
// 	ClassID  string `path:"class_id"`
// 	LessonID string `path:"lesson_id"`
// 	ExerciseID string `path:"exercise_id"`
// }

// Endpoint для обновления прогресса курса
type UpdateClassProgressInput struct {
	Status   string `json:"status" binding:"required"`
	CourseID string `path:"course_id"`
	ClassID  string `path:"class_id"`
}

type UpdateLessonProgressInput struct {
	Status   string `json:"status" binding:"required"`
	CourseID string `path:"course_id"`
	ClassID  string `path:"class_id"`
	LessonID string `path:"lesson_id"`
}

type UpdateExerciseProgressInput struct {
	Status     string `json:"status" binding:"required"`
	CourseID   string `path:"course_id"`
	ClassID    string `path:"class_id"`
	LessonID   string `path:"lesson_id"`
	ExerciseID string `path:"exercise_id"`
}

func GetFullClientProgress(c *gin.Context) (*course.ClientProgress, error) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	userClaims, err := auth.ExtractClaims(claims)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	progress, err := course.GetClientProgressByClientID(userClaims.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			newProgress := course.ClientProgress{
				ClientID:  userClaims.ID,
				Courses:   []course.CourseStatus{},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			progress, err = course.CreateClientProgress(&newProgress)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	database_course.EnsureFullClientStructure(userClaims.ID, progress)

	// Логирование прогресса перед возвратом (для отладки)
	// log.Println("Progress: ", progress)

	return progress, nil
}

func UpdateClassProgress(c *gin.Context, in *UpdateClassProgressInput) (*course.ClientProgress, error) {
	courseID, err := strconv.Atoi(in.CourseID)
	if err != nil {
		return nil, &gin.Error{
			Err:  errors.New("invalid course_id"),
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "invalid course_id"},
		}
	}

	classID, err := strconv.Atoi(in.ClassID)
	if err != nil {
		return nil, &gin.Error{
			Err:  errors.New("invalid class_id"),
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "invalid class_id"},
		}
	}

	claims := c.MustGet("claims").(jwt.MapClaims)
	userClaims, err := auth.ExtractClaims(claims)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	err = course.UpdateClassStatus(userClaims.ID, courseID, classID, in.Status)
	if err != nil {
		return nil, err
	}

	progress, err := course.GetClientProgressByClientAndCourseID(userClaims.ID, courseID)
	if err != nil {
		return nil, err
	}

	// Ensure the structure is fully populated for this specific course
	database_course.EnsureFullStructure(userClaims.ID, courseID, progress)

	return progress, nil
}

func UpdateLessonProgress(c *gin.Context, in *UpdateLessonProgressInput) (*course.ClientProgress, error) {
	classID, err := strconv.Atoi(in.ClassID)
	if err != nil {
		return nil, &gin.Error{
			Err:  errors.New("invalid class_id"),
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "invalid class_id"},
		}
	}

	lessonID, err := strconv.Atoi(in.LessonID)
	if err != nil {
		return nil, &gin.Error{
			Err:  errors.New("invalid lesson_id"),
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "invalid lesson_id"},
		}
	}

	claims := c.MustGet("claims").(jwt.MapClaims)
	userClaims, err := auth.ExtractClaims(claims)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	err = course.UpdateLessonStatus(userClaims.ID, classID, lessonID, in.Status)
	if err != nil {
		return nil, err
	}

	progress, err := course.GetClientProgressByClientAndCourseID(userClaims.ID, classID)
	if err != nil {
		return nil, err
	}

	// Ensure the structure is fully populated for this specific course
	database_course.EnsureFullStructure(userClaims.ID, classID, progress)

	return progress, nil
}

func UpdateExerciseProgress(c *gin.Context, in *UpdateExerciseProgressInput) (*course.ClientProgress, error) {
	lessonID, err := strconv.Atoi(in.LessonID)
	if err != nil {
		return nil, &gin.Error{
			Err:  errors.New("invalid lesson_id"),
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "invalid lesson_id"},
		}
	}

	exerciseID, err := strconv.Atoi(in.ExerciseID)
	if err != nil {
		return nil, &gin.Error{
			Err:  errors.New("invalid exercise_id"),
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "invalid exercise_id"},
		}
	}

	claims := c.MustGet("claims").(jwt.MapClaims)
	userClaims, err := auth.ExtractClaims(claims)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	err = course.UpdateExerciseStatus(userClaims.ID, lessonID, exerciseID, in.Status)
	if err != nil {
		return nil, err
	}

	progress, err := course.GetClientProgressByClientAndCourseID(userClaims.ID, lessonID)
	if err != nil {
		return nil, err
	}

	// Ensure the structure is fully populated for this specific course
	database_course.EnsureFullStructure(userClaims.ID, lessonID, progress)

	return progress, nil
}
