// internal/api/course/course.go
package course

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/niazlv/sport-plus-LCT/internal/api/auth"
	"github.com/niazlv/sport-plus-LCT/internal/database/course"
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
	api.GET("/:id", []fizz.OperationOption{fizz.Summary("Get course by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetCourseByID, 200))
	api.POST("", []fizz.OperationOption{fizz.Summary("Create a new course"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(CreateCourse, 201))
	api.PUT("/:id", []fizz.OperationOption{fizz.Summary("Update course by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(UpdateCourse, 200))
}

type CourseOutput struct {
	Course course.Course `json:"course"`
}

type CoursesOutput struct {
	Courses []course.Course `json:"courses"`
}
type GetCourseByIDParams struct {
	ID string `path:"id" binding:"required"`
}

func GetCourses(c *gin.Context) (*CoursesOutput, error) {

	var courses []course.Course
	result := db.Find(&courses)
	if result.Error != nil {
		return nil, result.Error
	}

	return &CoursesOutput{
		Courses: courses,
	}, nil
}

func GetCourseByID(c *gin.Context, params *GetCourseByIDParams) (*CourseOutput, error) {
	id := c.Param("id")

	var course course.Course
	result := db.First(&course, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, &gin.Error{
				Err:  result.Error,
				Type: gin.ErrorTypePublic,
				Meta: gin.H{"error": "course not found"},
			}
		}
		return nil, result.Error
	}

	return &CourseOutput{
		Course: course,
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
		return nil, result.Error
	}

	return &CourseOutput{
		Course: newCourse,
	}, nil
}

type UpdateCourseInput struct {
	ID                string  `path:"id" binding:"required"`
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
	id := c.Param("id")

	var course course.Course
	result := db.First(&course, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, &gin.Error{
				Err:  result.Error,
				Type: gin.ErrorTypePublic,
				Meta: gin.H{"error": "course not found"},
			}
		}
		return nil, result.Error
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

	result = db.Save(&course)
	if result.Error != nil {
		return nil, result.Error
	}

	return &CourseOutput{
		Course: course,
	}, nil
}
