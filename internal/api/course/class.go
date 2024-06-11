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

func SetupClassRoutes(api *fizz.RouterGroup) {
	classesAPI := api.Group("/:course_id/classes", "Classes", "Classes related endpoints")
	classesAPI.GET("", []fizz.OperationOption{fizz.Summary("Get list of classes for a course"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetClasses, 200))
	classesAPI.GET("/:class_id", []fizz.OperationOption{fizz.Summary("Get class by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetClassByID, 200))
	classesAPI.POST("", []fizz.OperationOption{fizz.Summary("Create a new class"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(CreateClass, 201))
	classesAPI.PUT("/:class_id", []fizz.OperationOption{fizz.Summary("Update class by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(UpdateClass, 200))
	classesAPI.DELETE("/:class_id", []fizz.OperationOption{fizz.Summary("Delete class by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(DeleteClass, 204))

	SetupClassImageRoutes(classesAPI)
}

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
	CourseID    string `path:"course_id" validate:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Cover       string `json:"cover"`
	Content     string `json:"content"`
	Video       string `json:"video"`
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
	Tips        string `json:"tips"`
}

type GetClassesParams struct {
	CourseID string `path:"course_id" binding:"required"`
}

func GetClasses(c *gin.Context, params *GetClassesParams) (*ClassesOutput, error) {
	courseID := params.CourseID
	log.Println("GetClasses called with course_id:", courseID)

	var classes []course.Class
	result := db.Preload("Images").Where("course_id = ?", courseID).Find(&classes)
	if result.Error != nil {
		log.Println("Error retrieving classes:", result.Error)
		return nil, result.Error
	}

	log.Printf("Retrieved classes for course_id %s: %+v\n", courseID, classes)
	return &ClassesOutput{
		Classes: classes,
	}, nil
}

func GetClassByID(c *gin.Context, params *GetClassByIDParams) (*ClassOutput, error) {
	classID := params.ID
	log.Println("GetClassByID called with class_id:", classID)

	var class course.Class
	result := db.Preload("Images").First(&class, classID)
	if result.Error != nil {
		log.Println("Error retrieving class:", result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			return nil, &gin.Error{
				Err:  result.Error,
				Type: gin.ErrorTypePublic,
				Meta: gin.H{"error": "class not found"},
			}
		}
		return nil, result.Error
	}

	log.Printf("Retrieved class: %+v\n", class)
	return &ClassOutput{
		Class: class,
	}, nil
}

func CreateClass(c *gin.Context, in *CreateClassInput) (*ClassOutput, error) {
	courseIDStr := in.CourseID
	log.Printf("CreateClass called with course_id: %s and input: %+v\n", courseIDStr, in)

	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		log.Println("Error converting course_id to int:", err)
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
		Tips:        in.Tips,
	}

	result := db.Create(&newClass)
	if result.Error != nil {
		log.Println("Error creating class:", result.Error)
		return nil, result.Error
	}

	log.Printf("Created class: %+v\n", newClass)
	return &ClassOutput{
		Class: newClass,
	}, nil
}

func UpdateClass(c *gin.Context, in *UpdateClassInput) (*ClassOutput, error) {
	classID := in.ID
	log.Printf("UpdateClass called with class_id: %s and input: %+v\n", classID, in)

	var class course.Class
	result := db.First(&class, classID)
	if result.Error != nil {
		log.Println("Error retrieving class:", result.Error)
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
	if in.Tips != "" {
		class.Tips = in.Tips
	}

	result = db.Save(&class)
	if result.Error != nil {
		log.Println("Error updating class:", result.Error)
		return nil, result.Error
	}

	log.Printf("Updated class: %+v\n", class)
	return &ClassOutput{
		Class: class,
	}, nil
}

func DeleteClass(c *gin.Context) error {
	classID := c.Param("class_id")
	log.Println("DeleteClass called with class_id:", classID)

	result := db.Delete(&course.Class{}, classID)
	if result.Error != nil {
		log.Println("Error deleting class:", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Println("Class not found with ID:", classID)
		return &gin.Error{
			Err:  gorm.ErrRecordNotFound,
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "class not found"},
		}
	}

	log.Printf("Deleted class with ID: %s\n", classID)
	return nil
}
