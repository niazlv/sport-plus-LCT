// internal/api/course/course.go
package course

import (
	"log"
	"strconv"

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
	api.GET("/:course_id", []fizz.OperationOption{fizz.Summary("Get course by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetCourseByID, 200))
	api.POST("", []fizz.OperationOption{fizz.Summary("Create a new course"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(CreateCourse, 201))
	api.PUT("/:course_id", []fizz.OperationOption{fizz.Summary("Update course by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(UpdateCourse, 200))
	// api.DELETE("/:course_id", []fizz.OperationOption{fizz.Summary("Delete course by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(DeleteCourse, 204))

	classesAPI := api.Group("/:course_id/classes", "Classes", "Classes related endpoints")
	classesAPI.GET("", []fizz.OperationOption{fizz.Summary("Get list of classes for a course"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetClasses, 200))
	classesAPI.GET("/:class_id", []fizz.OperationOption{fizz.Summary("Get class by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetClassByID, 200))
	classesAPI.POST("", []fizz.OperationOption{fizz.Summary("Create a new class"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(CreateClass, 201))
	classesAPI.PUT("/:class_id", []fizz.OperationOption{fizz.Summary("Update class by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(UpdateClass, 200))
	classesAPI.DELETE("/:class_id", []fizz.OperationOption{fizz.Summary("Delete class by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(DeleteClass, 204))
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
	id := c.Param("course_id")

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
	ID                string  `path:"course_id" binding:"required"`
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
	id := c.Param("course_id")

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

// Endpoints for Class

type ClassOutput struct {
	Class course.Class `json:"class"`
}

type ClassesOutput struct {
	Classes []course.Class `json:"classes"`
}

type GetClassByIDParams struct {
	CourseID string `path:"course_id" binding:"required"`
	ID       string `path:"class_id" binding:"required"`
}

type CreateClassInput struct {
	CourseID    string `path:"course_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Cover       string `json:"cover"`
	Content     string `json:"content"`
	Video       string `json:"video"`
	Image       string `json:"image"`
	Tips        string `json:"tips"`
}

type UpdateClassInput struct {
	CourseID    string `path:"course_id" binding:"required"`
	ID          string `path:"class_id" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Cover       string `json:"cover"`
	Content     string `json:"content"`
	Video       string `json:"video"`
	Image       string `json:"image"`
	Tips        string `json:"tips"`
}

func GetClasses(c *gin.Context) (*ClassesOutput, error) {
	courseID := c.Param("course_id")

	var classes []course.Class
	result := db.Where("course_id = ?", courseID).Find(&classes)
	if result.Error != nil {
		return nil, result.Error
	}

	return &ClassesOutput{
		Classes: classes,
	}, nil
}

func GetClassByID(c *gin.Context, params *GetClassByIDParams) (*ClassOutput, error) {
	classID := params.ID

	var class course.Class
	result := db.First(&class, classID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, &gin.Error{
				Err:  result.Error,
				Type: gin.ErrorTypePublic,
				Meta: gin.H{"error": "class not found"},
			}
		}
		return nil, result.Error
	}

	return &ClassOutput{
		Class: class,
	}, nil
}

func CreateClass(c *gin.Context, in *CreateClassInput) (*ClassOutput, error) {
	courseIDStr := c.Param("course_id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		return nil, &gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "invalid course_id"},
		}
	}

	newClass := course.Class{
		CourseID:    courseID,
		Title:       in.Title,
		Description: in.Description,
		Cover:       in.Cover,
		Content:     in.Content,
		Video:       in.Video,
		Image:       in.Image,
		Tips:        in.Tips,
	}

	result := db.Create(&newClass)
	if result.Error != nil {
		return nil, result.Error
	}

	return &ClassOutput{
		Class: newClass,
	}, nil
}

func UpdateClass(c *gin.Context, in *UpdateClassInput) (*ClassOutput, error) {
	classID := in.ID

	var class course.Class
	result := db.First(&class, classID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, &gin.Error{
				Err:  result.Error,
				Type: gin.ErrorTypePublic,
				Meta: gin.H{"error": "class not found"},
			}
		}
		return nil, result.Error
	}

	if in.Title != "" {
		class.Title = in.Title
	}
	if in.Description != "" {
		class.Description = in.Description
	}
	if in.Cover != "" {
		class.Cover = in.Cover
	}
	if in.Content != "" {
		class.Content = in.Content
	}
	if in.Video != "" {
		class.Video = in.Video
	}
	if in.Image != "" {
		class.Image = in.Image
	}
	if in.Tips != "" {
		class.Tips = in.Tips
	}

	result = db.Save(&class)
	if result.Error != nil {
		return nil, result.Error
	}

	return &ClassOutput{
		Class: class,
	}, nil
}

func DeleteClass(c *gin.Context) error {
	classID := c.Param("id")

	result := db.Delete(&course.Class{}, classID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return &gin.Error{
			Err:  gorm.ErrRecordNotFound,
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "class not found"},
		}
	}
	return nil
}
