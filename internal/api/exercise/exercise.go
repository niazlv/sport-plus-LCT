package exercise

import (
	"errors"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/niazlv/sport-plus-LCT/internal/api/auth"
	"github.com/niazlv/sport-plus-LCT/internal/database/exercise"
	exercise_class "github.com/niazlv/sport-plus-LCT/internal/database/exercise"
	"github.com/wI2L/fizz"
	"gorm.io/gorm"
)

var db *gorm.DB

func Setup(rg *fizz.RouterGroup) {
	api := rg.Group("exercise", "Exercise", "Exercise related endpoints")

	var err error
	db, err = exercise.InitDB()
	if err != nil {
		log.Fatal("db exercises can't be init: ", err)
	}

	api.GET("", []fizz.OperationOption{fizz.Summary("Get list of exercises"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetExercises, 200))
	api.GET("/:exercise_id", []fizz.OperationOption{fizz.Summary("Get exercise by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetExerciseByID, 200))
	api.GET("/filter", []fizz.OperationOption{fizz.Summary("Filter exercises"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(FilterExercises, 200))
	api.POST("", []fizz.OperationOption{fizz.Summary("Create a new exercise"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(CreateExercise, 201))
	api.PUT("/:exercise_id", []fizz.OperationOption{fizz.Summary("Update exercise by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(UpdateExercise, 200))
	api.DELETE("/:exercise_id", []fizz.OperationOption{fizz.Summary("Delete exercise by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(DeleteExercise, 204))
}

type ExerciseOutput struct {
	Exercise exercise.Exercise `json:"exercise"`
}

type ExercisesOutput struct {
	Exercises []exercise.Exercise `json:"exercises"`
}

type GetExerciseByIDParams struct {
	ID string `path:"exercise_id" binding:"required"`
}

type FilterExercisesParams struct {
	AdditionalMuscle string `query:"additional_muscle"`
	Muscle           string `query:"muscle"`
	Difficulty       string `query:"difficulty"`
}

func GetExercises(c *gin.Context) (*ExercisesOutput, error) {
	log.Println("GetExercises called")

	var exercises []exercise.Exercise
	result := db.Preload("Photos").Find(&exercises)
	if result.Error != nil {
		log.Println("Error retrieving exercises:", result.Error)
		return nil, result.Error
	}

	log.Printf("Retrieved exercises: %+v\n", exercises)
	return &ExercisesOutput{
		Exercises: exercises,
	}, nil
}

func GetExerciseByID(c *gin.Context, params *GetExerciseByIDParams) (*ExerciseOutput, error) {
	idStr := params.ID
	log.Println("GetExerciseByID called with ID:", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid exercise ID:", idStr)
		return nil, gin.Error{
			Err:  errors.New("invalid exercise_id"),
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "invalid exercise_id"},
		}
	}

	var exercise exercise.Exercise
	result := db.Preload("Photos").First(&exercise, id)
	if result.Error != nil {
		log.Println("Error retrieving exercise:", result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			return nil, gin.Error{
				Err:  result.Error,
				Type: gin.ErrorTypePublic,
				Meta: gin.H{"error": "exercise not found"},
			}
		}
		return nil, result.Error
	}

	log.Printf("Retrieved exercise: %+v\n", exercise)
	return &ExerciseOutput{
		Exercise: exercise,
	}, nil
}

func FilterExercises(c *gin.Context, params *FilterExercisesParams) (*ExercisesOutput, error) {
	log.Println("FilterExercises called with params:", params)

	query := db.Preload("Photos")
	if params.AdditionalMuscle != "" {
		query = query.Where("additional_muscle LIKE ?", "%"+params.AdditionalMuscle+"%")
	}
	if params.Muscle != "" {
		query = query.Where("muscle LIKE ?", "%"+params.Muscle+"%")
	}
	if params.Difficulty != "" {
		query = query.Where("difficulty = ?", params.Difficulty)
	}

	var exercises []exercise.Exercise
	result := query.Find(&exercises)
	if result.Error != nil {
		log.Println("Error retrieving exercises:", result.Error)
		return nil, result.Error
	}

	log.Printf("Filtered exercises: %+v\n", exercises)
	return &ExercisesOutput{
		Exercises: exercises,
	}, nil
}

type CreateExerciseInput struct {
	OriginalUri      string   `json:"original_uri" binding:"required"`
	Name             string   `json:"name" binding:"required"`
	Muscle           string   `json:"muscle" binding:"required"`
	AdditionalMuscle string   `json:"additional_muscle" binding:"required"`
	Type             string   `json:"type" binding:"required"`
	Equipment        string   `json:"equipment" binding:"required"`
	Difficulty       string   `json:"difficulty" binding:"required"`
	Photos           []string `json:"photos"`
}

func CreateExercise(c *gin.Context, in *CreateExerciseInput) (*ExerciseOutput, error) {
	log.Printf("CreateExercise called with input: %+v\n", in)

	photos := make([]exercise.Photo, len(in.Photos))
	for i, url := range in.Photos {
		photos[i] = exercise.Photo{URL: url}
	}

	newExercise := exercise.Exercise{
		OriginalUri:      in.OriginalUri,
		Name:             in.Name,
		Muscle:           in.Muscle,
		AdditionalMuscle: in.AdditionalMuscle,
		Type:             in.Type,
		Equipment:        in.Equipment,
		Difficulty:       in.Difficulty,
		Photos:           photos,
	}

	result := db.Create(&newExercise)
	if result.Error != nil {
		log.Println("Error creating exercise:", result.Error)
		return nil, result.Error
	}

	log.Printf("Created exercise: %+v\n", newExercise)
	return &ExerciseOutput{
		Exercise: newExercise,
	}, nil
}

type UpdateExerciseInput struct {
	ID               string   `path:"exercise_id" binding:"required"`
	OriginalUri      string   `json:"original_uri"`
	Name             string   `json:"name"`
	Muscle           string   `json:"muscle"`
	AdditionalMuscle string   `json:"additional_muscle"`
	Type             string   `json:"type"`
	Equipment        string   `json:"equipment"`
	Difficulty       string   `json:"difficulty"`
	Photos           []string `json:"photos"`
}

func UpdateExercise(c *gin.Context, in *UpdateExerciseInput) (*ExerciseOutput, error) {
	id := in.ID
	log.Printf("UpdateExercise called with ID: %s and input: %+v\n", id, in)

	var exercise exercise.Exercise
	result := db.First(&exercise, id)
	if result.Error != nil {
		log.Println("Error retrieving exercise:", result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			return nil, gin.Error{
				Err:  result.Error,
				Type: gin.ErrorTypePublic,
				Meta: gin.H{"error": "exercise not found"},
			}
		}
		return nil, result.Error
	}

	if in.OriginalUri != "" {
		exercise.OriginalUri = in.OriginalUri
	}
	if in.Name != "" {
		exercise.Name = in.Name
	}
	if in.Muscle != "" {
		exercise.Muscle = in.Muscle
	}
	if in.AdditionalMuscle != "" {
		exercise.AdditionalMuscle = in.AdditionalMuscle
	}
	if in.Type != "" {
		exercise.Type = in.Type
	}
	if in.Equipment != "" {
		exercise.Equipment = in.Equipment
	}
	if in.Difficulty != "" {
		exercise.Difficulty = in.Difficulty
	}
	if len(in.Photos) > 0 {
		photos := make([]exercise_class.Photo, len(in.Photos))
		for i, url := range in.Photos {
			photos[i] = exercise_class.Photo{URL: url, ExerciseID: exercise.Id}
		}
		exercise.Photos = photos
	}

	result = db.Save(&exercise)
	if result.Error != nil {
		log.Println("Error updating exercise:", result.Error)
		return nil, result.Error
	}

	log.Printf("Updated exercise: %+v\n", exercise)
	return &ExerciseOutput{
		Exercise: exercise,
	}, nil
}

type DeleteExerciseParams struct {
	ID string `path:"exercise_id" binding:"required"`
}

func DeleteExercise(c *gin.Context, params *DeleteExerciseParams) error {
	idStr := params.ID
	log.Println("DeleteExercise called with ID:", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid exercise ID:", idStr)
		return gin.Error{
			Err:  errors.New("invalid exercise_id"),
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "invalid exercise_id"},
		}
	}

	result := db.Delete(&exercise.Exercise{}, id)
	if result.Error != nil {
		log.Println("Error deleting exercise:", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gin.Error{
			Err:  gorm.ErrRecordNotFound,
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "exercise not found"},
		}
	}

	log.Printf("Deleted exercise with ID: %d\n", id)
	return nil
}
