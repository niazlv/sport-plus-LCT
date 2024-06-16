package review

import (
	"errors"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/niazlv/sport-plus-LCT/internal/api/auth"
	"github.com/niazlv/sport-plus-LCT/internal/database/review"
	database_review "github.com/niazlv/sport-plus-LCT/internal/database/review"
	"github.com/wI2L/fizz"
	"gorm.io/gorm"
)

var db *gorm.DB

func Setup(rg *fizz.RouterGroup) {
	api := rg.Group("review", "Review", "Review related endpoints")

	var err error
	db, err = review.InitDB()
	if err != nil {
		log.Fatal("db reviews can't be init: ", err)
	}

	api.GET("/:review_id", []fizz.OperationOption{fizz.Summary("Get review by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetReviewByID, 200))
	api.GET("/class/:class_id", []fizz.OperationOption{fizz.Summary("Get reviews by Class ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetReviewsByClassID, 200))
	api.POST("", []fizz.OperationOption{fizz.Summary("Create a new review"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(CreateReview, 201))
	api.PUT("/:review_id", []fizz.OperationOption{fizz.Summary("Update review by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(UpdateReview, 200))
	api.DELETE("/:review_id", []fizz.OperationOption{fizz.Summary("Delete review by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(DeleteReview, 204))
}

type ReviewOutput struct {
	Review review.Review `json:"review"`
}

type ReviewsOutput struct {
	Reviews []review.Review `json:"reviews"`
}

type GetReviewByIDParams struct {
	ID string `path:"review_id" binding:"required"`
}

type GetReviewsByClassIDParams struct {
	ClassID string `path:"class_id" binding:"required"`
}

func GetReviewByID(c *gin.Context, params *GetReviewByIDParams) (*ReviewOutput, error) {
	idStr := params.ID
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, gin.Error{
			Err:  errors.New("invalid review_id"),
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "invalid review_id"},
		}
	}

	review, err := review.GetReviewByID(id)
	if err != nil {
		return nil, err
	}

	return &ReviewOutput{
		Review: *review,
	}, nil
}

func GetReviewsByClassID(c *gin.Context, params *GetReviewsByClassIDParams) (*ReviewsOutput, error) {
	classIDStr := params.ClassID
	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		return nil, gin.Error{
			Err:  errors.New("invalid class_id"),
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "invalid class_id"},
		}
	}

	reviews, err := review.GetReviewsByClassID(classID)
	if err != nil {
		return nil, err
	}

	return &ReviewsOutput{
		Reviews: reviews,
	}, nil
}

type CreateReviewInput struct {
	ClassID int `json:"class_id"`
	//ClientID         int    `json:"client_id" binding:"required"`
	TrainerID        int    `json:"trainer_id"`
	DifficultyRating int    `json:"difficulty_rating"`
	WellBeingRating  int    `json:"well_being_rating"`
	OverallRating    int    `json:"overall_rating"`
	Comment          string `json:"comment"`
}

func CreateReview(c *gin.Context, in *CreateReviewInput) (*ReviewOutput, error) {
	claims := c.MustGet("claims").(jwt.MapClaims)

	userClaims, err := auth.ExtractClaims(claims)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	newReview := review.Review{
		ClassID:          in.ClassID,
		ClientID:         userClaims.ID,
		TrainerID:        in.TrainerID,
		DifficultyRating: in.DifficultyRating,
		WellBeingRating:  in.WellBeingRating,
		OverallRating:    in.OverallRating,
		Comment:          in.Comment,
	}

	result, err := review.CreateReview(&newReview)
	if err != nil {
		return nil, err
	}

	return &ReviewOutput{
		Review: *result,
	}, nil
}

type UpdateReviewInput struct {
	ID               string `path:"review_id" binding:"required"`
	DifficultyRating int    `json:"difficulty_rating"`
	WellBeingRating  int    `json:"well_being_rating"`
	OverallRating    int    `json:"overall_rating"`
	Comment          string `json:"comment"`
}

func UpdateReview(c *gin.Context, in *UpdateReviewInput) (*ReviewOutput, error) {
	idStr := in.ID
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, gin.Error{
			Err:  errors.New("invalid review_id"),
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "invalid review_id"},
		}
	}

	review, err := review.GetReviewByID(id)
	if err != nil {
		return nil, err
	}

	if in.DifficultyRating != 0 {
		review.DifficultyRating = in.DifficultyRating
	}
	if in.WellBeingRating != 0 {
		review.WellBeingRating = in.WellBeingRating
	}
	if in.OverallRating != 0 {
		review.OverallRating = in.OverallRating
	}
	if in.Comment != "" {
		review.Comment = in.Comment
	}

	err = database_review.UpdateReview(review)
	if err != nil {
		return nil, err
	}

	return &ReviewOutput{
		Review: *review,
	}, nil
}

type DeleteReviewParams struct {
	ID string `path:"review_id" binding:"required"`
}

func DeleteReview(c *gin.Context, params *DeleteReviewParams) error {
	idStr := params.ID
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return gin.Error{
			Err:  errors.New("invalid review_id"),
			Type: gin.ErrorTypePublic,
			Meta: gin.H{"error": "invalid review_id"},
		}
	}

	err = review.DeleteReview(id)
	if err != nil {
		return err
	}

	return nil
}
