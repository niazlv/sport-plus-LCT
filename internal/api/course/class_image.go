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

func SetupClassImageRoutes(api *fizz.RouterGroup) {
	imagesAPI := api.Group("/:class_id/images", "Images", "Images related endpoints")
	imagesAPI.POST("", []fizz.OperationOption{fizz.Summary("Create a new image for class"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(CreateClassImage, 201))
	imagesAPI.GET("/:image_id", []fizz.OperationOption{fizz.Summary("Get image by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetClassImageByID, 200))
	imagesAPI.PUT("/:image_id", []fizz.OperationOption{fizz.Summary("Update image by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(UpdateClassImage, 200))
	imagesAPI.DELETE("/:image_id", []fizz.OperationOption{fizz.Summary("Delete image by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(DeleteClassImage, 204))
}

type ClassImageOutput struct {
	ClassImage course.ClassImage `json:"class_image"`
}

type ClassImagesOutput struct {
	ClassImages []course.ClassImage `json:"class_images"`
}

type GetClassImageByIDParams struct {
	CourseID string `path:"course_id" binding:"required"`
	ClassID  string `path:"class_id" binding:"required"`
	ID       string `path:"image_id" binding:"required"`
}

type CreateClassImageInput struct {
	CourseID string `path:"course_id" validate:"required"`
	ClassID  string `path:"class_id" validate:"required"`
	Image    string `json:"image" binding:"required"`
}

type UpdateClassImageInput struct {
	CourseID string `path:"course_id" binding:"required"`
	ClassID  string `path:"class_id" binding:"required"`
	ID       string `path:"image_id" binding:"required"`
	Image    string `json:"image"`
}

type GetClassImagesParams struct {
	CourseID string `path:"course_id" binding:"required"`
	ClassID  string `path:"class_id" binding:"required"`
}

func CreateClassImage(c *gin.Context, in *CreateClassImageInput) (*ClassImageOutput, error) {
	classIDStr := in.ClassID
	log.Printf("CreateClassImage called with class_id: %s and input: %+v\n", classIDStr, in)

	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		log.Println("Error converting class_id to int:", err)
		return nil, &gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "invalid class_id"},
		}
	}

	newClassImage := course.ClassImage{
		ClassID: classID,
		Image:   in.Image,
	}

	result := db.Create(&newClassImage)
	if result.Error != nil {
		log.Println("Error creating class image:", result.Error)
		return nil, result.Error
	}

	log.Printf("Created class image: %+v\n", newClassImage)
	return &ClassImageOutput{
		ClassImage: newClassImage,
	}, nil
}

func GetClassImageByID(c *gin.Context, params *GetClassImageByIDParams) (*ClassImageOutput, error) {
	imageID := params.ID
	log.Println("GetClassImageByID called with image_id:", imageID)

	var classImage course.ClassImage
	result := db.First(&classImage, imageID)
	if result.Error != nil {
		log.Println("Error retrieving class image:", result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			return nil, &gin.Error{
				Err:  result.Error,
				Type: gin.ErrorTypePublic,
				Meta: gin.H{"error": "class image not found"},
			}
		}
		return nil, result.Error
	}

	log.Printf("Retrieved class image: %+v\n", classImage)
	return &ClassImageOutput{
		ClassImage: classImage,
	}, nil
}

func UpdateClassImage(c *gin.Context, in *UpdateClassImageInput) (*ClassImageOutput, error) {
	imageID := in.ID
	log.Printf("UpdateClassImage called with image_id: %s and input: %+v\n", imageID, in)

	var classImage course.ClassImage
	result := db.First(&classImage, imageID)
	if result.Error != nil {
		log.Println("Error retrieving class image:", result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			return nil, &gin.Error{
				Err:  result.Error,
				Type: gin.ErrorTypePublic,
				Meta: gin.H{"error": "class image not found"},
			}
		}
		return nil, result.Error
	}

	if in.Image != "" {
		classImage.Image = in.Image
	}

	result = db.Save(&classImage)
	if result.Error != nil {
		log.Println("Error updating class image:", result.Error)
		return nil, result.Error
	}

	log.Printf("Updated class image: %+v\n", classImage)
	return &ClassImageOutput{
		ClassImage: classImage,
	}, nil
}

func DeleteClassImage(c *gin.Context) error {
	imageID := c.Param("image_id")
	log.Println("DeleteClassImage called with image_id:", imageID)

	result := db.Delete(&course.ClassImage{}, imageID)
	if result.Error != nil {
		log.Println("Error deleting class image:", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Println("Class image not found with ID:", imageID)
		return &gin.Error{
			Err:  gorm.ErrRecordNotFound,
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "class image not found"},
		}
	}

	log.Printf("Deleted class image with ID: %s\n", imageID)
	return nil
}
