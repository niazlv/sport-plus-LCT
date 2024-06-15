// internal/api/course/lesson.go

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

func SetupLessonRoutes(api *fizz.RouterGroup) {
	lessonsAPI := api.Group("/:class_id/lessons", "Lessons", "Lessons related endpoints")
	lessonsAPI.GET("", []fizz.OperationOption{fizz.Summary("Get list of lessons for a class"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetLessons, 200))
	lessonsAPI.GET("/:lesson_id", []fizz.OperationOption{fizz.Summary("Get lesson by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetLessonByID, 200))
	lessonsAPI.POST("", []fizz.OperationOption{fizz.Summary("Create a new lesson"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(CreateLesson, 201))
	lessonsAPI.PUT("/:lesson_id", []fizz.OperationOption{fizz.Summary("Update lesson by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(UpdateLesson, 200))
	lessonsAPI.DELETE("/:lesson_id", []fizz.OperationOption{fizz.Summary("Delete lesson by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(DeleteLesson, 204))

	SetupClassImageRoutes(lessonsAPI)
}

type LessonOutput struct {
	Lesson course.Lesson `json:"lesson"`
}

type LessonsOutput struct {
	Lessons []course.Lesson `json:"lessons"`
}

type GetLessonByIDParams struct {
	CourseID string `path:"course_id" binding:"required"`
	ClassID  string `path:"class_id" binding:"required"`
	ID       string `path:"lesson_id" binding:"required"`
}

type CreateLessonInput struct {
	CourseID        string `path:"course_id" validate:"required"`
	ClassID         string `path:"class_id" validate:"required"`
	Title           string `json:"title" binding:"required"`
	Description     string `json:"description"`
	Video           string `json:"video"`
	Tips            string `json:"tips"`
	DurationSeconds int    `json:"duration_seconds"`
}

type UpdateLessonInput struct {
	CourseID        string `path:"course_id" binding:"required"`
	ClassID         string `path:"class_id" binding:"required"`
	ID              string `path:"lesson_id" binding:"required"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Video           string `json:"video"`
	Tips            string `json:"tips"`
	DurationSeconds int    `json:"duration_seconds"`
}

type GetLessonsParams struct {
	CourseID string `path:"course_id" binding:"required"`
	ClassID  string `path:"class_id" binding:"required"`
}

func GetLessons(c *gin.Context, params *GetLessonsParams) (*LessonsOutput, error) {
	classID := params.ClassID
	log.Println("GetLessons called with class_id:", classID)

	var lessons []course.Lesson
	result := db.Preload("Images").Where("class_id = ?", classID).Find(&lessons)
	if result.Error != nil {
		log.Println("Error retrieving lessons:", result.Error)
		return nil, result.Error
	}

	log.Printf("Retrieved lessons for class_id %s: %+v\n", classID, lessons)
	return &LessonsOutput{
		Lessons: lessons,
	}, nil
}

func GetLessonByID(c *gin.Context, params *GetLessonByIDParams) (*LessonOutput, error) {
	lessonID := params.ID
	log.Println("GetLessonByID called with lesson_id:", lessonID)

	var lesson course.Lesson
	result := db.Preload("Images").First(&lesson, lessonID)
	if result.Error != nil {
		log.Println("Error retrieving lesson:", result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			return nil, &gin.Error{
				Err:  result.Error,
				Type: gin.ErrorTypePublic,
				Meta: gin.H{"error": "lesson not found"},
			}
		}
		return nil, result.Error
	}

	log.Printf("Retrieved lesson: %+v\n", lesson)
	return &LessonOutput{
		Lesson: lesson,
	}, nil
}

func CreateLesson(c *gin.Context, in *CreateLessonInput) (*LessonOutput, error) {
	classIDStr := in.ClassID
	log.Printf("CreateLesson called with class_id: %s and input: %+v\n", classIDStr, in)

	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		log.Println("Error converting class_id to int:", err)
		return nil, &gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "invalid class_id"},
		}
	}

	courseIDStr := in.CourseID
	log.Printf("CreateLesson called with course_id: %s and input: %+v\n", courseIDStr, in)

	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		log.Println("Error converting course_id to int:", err)
		return nil, &gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "invalid course_id"},
		}
	}

	newLesson := course.Lesson{
		CourseID:        courseID,
		ClassID:         classID,
		Title:           in.Title,
		Description:     in.Description,
		Video:           in.Video,
		Tips:            in.Tips,
		DurationSeconds: in.DurationSeconds,
	}

	result := db.Create(&newLesson)
	if result.Error != nil {
		log.Println("Error creating lesson:", result.Error)
		return nil, result.Error
	}

	log.Printf("Created lesson: %+v\n", newLesson)
	return &LessonOutput{
		Lesson: newLesson,
	}, nil
}

func UpdateLesson(c *gin.Context, in *UpdateLessonInput) (*LessonOutput, error) {
	lessonID := in.ID
	log.Printf("UpdateLesson called with lesson_id: %s and input: %+v\n", lessonID, in)

	var lesson course.Lesson
	result := db.First(&lesson, lessonID)
	if result.Error != nil {
		log.Println("Error retrieving lesson:", result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			return nil, &gin.Error{
				Err:  result.Error,
				Type: gin.ErrorTypePublic,
				Meta: gin.H{"error": "lesson not found"},
			}
		}
		return nil, result.Error
	}

	if in.Title != "" {
		lesson.Title = in.Title
	}
	if in.Description != "" {
		lesson.Description = in.Description
	}
	if in.Video != "" {
		lesson.Video = in.Video
	}
	if in.Tips != "" {
		lesson.Tips = in.Tips
	}
	if in.DurationSeconds != 0 {
		lesson.DurationSeconds = in.DurationSeconds
	}

	result = db.Save(&lesson)
	if result.Error != nil {
		log.Println("Error updating lesson:", result.Error)
		return nil, result.Error
	}

	log.Printf("Updated lesson: %+v\n", lesson)
	return &LessonOutput{
		Lesson: lesson,
	}, nil
}

func DeleteLesson(c *gin.Context) error {
	lessonID := c.Param("lesson_id")
	log.Println("DeleteLesson called with lesson_id:", lessonID)

	result := db.Delete(&course.Lesson{}, lessonID)
	if result.Error != nil {
		log.Println("Error deleting lesson:", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Println("Lesson not found with ID:", lessonID)
		return &gin.Error{
			Err:  gorm.ErrRecordNotFound,
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "lesson not found"},
		}
	}

	log.Printf("Deleted lesson with ID: %s\n", lessonID)
	return nil
}
