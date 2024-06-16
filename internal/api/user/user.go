package user

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/juju/errors"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/niazlv/sport-plus-LCT/internal/api/auth"
	database "github.com/niazlv/sport-plus-LCT/internal/database/auth"
	"github.com/wI2L/fizz"
	"gorm.io/gorm"
)

type JWTAuthSecurityType struct {
	Name         string
	Type         string
	BearerFormat string
	Scheme       string
}

func Setup(rg *fizz.RouterGroup) {
	api := rg.Group("user", "User", "User related endpoints")

	_ = api
	api.GET("", []fizz.OperationOption{fizz.Summary("Return Your User"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetUser, 200))
	api.GET("/:id", []fizz.OperationOption{fizz.Summary("Return User by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetUserByID, 200))
	api.PUT("/onboarding", []fizz.OperationOption{fizz.Summary("Update User data after onboarding"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(putOnboarding, 200))
	api.POST("/upload/icon", []fizz.OperationOption{fizz.Summary("Upload user icon"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(UploadUserIcon, 201))

	api.POST("/measurements", []fizz.OperationOption{fizz.Summary("Add a new measurement"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(AddMeasurement, 201))
	api.PUT("/measurements/:id", []fizz.OperationOption{fizz.Summary("Update a measurement"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(UpdateMeasurement, 200))
	api.DELETE("/measurements/:id", []fizz.OperationOption{fizz.Summary("Delete a measurement"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(DeleteMeasurement, 204))

	api.POST("/trains", []fizz.OperationOption{fizz.Summary("Add a new train"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(AddTrain, 201))
	api.PUT("/trains/:id", []fizz.OperationOption{fizz.Summary("Update a train"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(UpdateTrain, 200))
	api.DELETE("/trains/:id", []fizz.OperationOption{fizz.Summary("Delete a train"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(DeleteTrain, 204))
}

type GetUserOutput struct {
	User database.User `json:"user"`
}

func GetUser(c *gin.Context) (*GetUserOutput, error) {
	claims := c.MustGet("claims").(jwt.MapClaims)

	userClaims, err := auth.ExtractClaims(claims)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	User, err := database.FindUserByID(userClaims.ID)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &GetUserOutput{
		User: *User,
	}, nil
}

type GetUserByIDInput struct {
	ID int `json:"id" path:"id" validate:"required" binding:"required"`
}

func GetUserByID(c *gin.Context, in *GetUserByIDInput) (*GetUserOutput, error) {
	User, err := database.FindUserByID(in.ID)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &GetUserOutput{
		User: *User,
	}, nil
}

type putOnboardingOutput struct {
	Status string `json:"status"`
}

func putOnboarding(c *gin.Context, in *database.User) (*putOnboardingOutput, error) {
	// get data from token
	claims := c.MustGet("claims").(jwt.MapClaims)
	userClaims, err := auth.ExtractClaims(claims)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Преобразуем входные данные в структуру User
	user := database.User{
		Id:               userClaims.ID,
		Gender:           in.Gender,
		Height:           in.Height,
		Weight:           in.Weight,
		Goals:            in.Goals,
		Experience:       in.Experience,
		GymMember:        in.GymMember,
		Beginner:         in.Beginner,
		GymName:          in.GymName,
		HealthConditions: in.HealthConditions,
		Role:             in.Role,
		Name:             in.Name,
		Icon:             in.Icon,
		About:            in.About,
		Achivements:      in.Achivements,
		Age:              in.Age,
	}

	// Обновляем пользователя в базе данных
	if err := database.PartialUpdateUser(&user); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gin.Error{
				Err:  err,
				Type: gin.ErrorTypePublic,
				Meta: gin.H{"error": "user not found"},
			}
		}
		return nil, gin.Error{
			Err:  err,
			Type: gin.ErrorTypePrivate,
			Meta: gin.H{"error": err.Error()},
		}
	}

	// Сохраняем измерения роста и веса
	if err := database.SaveMeasurements(in.Height, user.Id, database.TypeHeight); err != nil {
		return nil, err
	}
	if err := database.SaveMeasurements(in.Weight, user.Id, database.TypeWeight); err != nil {
		return nil, err
	}
	if err := database.SaveMeasurements(in.Water, user.Id, database.TypeWater); err != nil {
		return nil, err
	}

	return &putOnboardingOutput{Status: "user updated successfully"}, nil
}

type UploadUserIconOutput struct {
	Url string `json:"url"`
}

func UploadUserIcon(c *gin.Context) (*UploadUserIconOutput, error) {
	claims := c.MustGet("claims").(jwt.MapClaims)

	userClaims, err := auth.ExtractClaims(claims)
	if err != nil {
		return nil, err
	}

	file, err := c.FormFile("icon")
	if err != nil {
		return nil, err
	}

	uploadDir := "./uploads/icons"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, err
	}

	filePath := filepath.Join(uploadDir, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return nil, err
	}

	iconURL := fmt.Sprintf("http://%s/uploads/icons/%s", c.Request.Host, url.PathEscape(file.Filename))

	// Обновляем иконку пользователя в базе данных
	user, err := database.FindUserByID(userClaims.ID)
	if err != nil {
		return nil, err
	}

	user.Icon = iconURL
	if err := database.PartialUpdateUser(user); err != nil {
		return nil, err
	}

	return &UploadUserIconOutput{
		Url: iconURL,
	}, nil
}

type AddMeasurementInput struct {
	// UserID int    `json:"userId" binding:"required"`
	Date  string `json:"date" binding:"required"`
	Value string `json:"value" binding:"required"`
	Type  string `json:"type" binding:"required"`
}

func (input *AddMeasurementInput) Validate() error {
	for _, validType := range database.ValidTypesMeasurement {
		if input.Type == validType {
			return nil
		}
	}
	return fmt.Errorf("invalid type: %s", input.Type)
}

type AddMeasurementOutput struct {
	Measurement database.Measurement `json:"measurement"`
}

func AddMeasurement(c *gin.Context, in *AddMeasurementInput) (*AddMeasurementOutput, error) {
	claims := c.MustGet("claims").(jwt.MapClaims)

	// Validate input
	if err := in.Validate(); err != nil {
		return nil, err
	}

	userClaims, err := auth.ExtractClaims(claims)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	User, err := database.FindUserByID(userClaims.ID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	measurement := &database.Measurement{
		UserID: User.Id,
		Date:   in.Date,
		Value:  in.Value,
		Type:   in.Type,
	}
	createdMeasurement, err := database.AddMeasurement(measurement)
	if err != nil {
		return nil, err
	}
	return &AddMeasurementOutput{
		Measurement: *createdMeasurement,
	}, nil
}

type UpdateMeasurementInput struct {
	ID int `json:"id" path:"id" binding:"required"`
	// UserID int    `json:"userId" binding:"required"`
	Date  string `json:"date" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type UpdateMeasurementOutput struct {
	Measurement database.Measurement `json:"measurement"`
}

func UpdateMeasurement(c *gin.Context, in *UpdateMeasurementInput) (*UpdateMeasurementOutput, error) {
	claims := c.MustGet("claims").(jwt.MapClaims)

	userClaims, err := auth.ExtractClaims(claims)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	User, err := database.FindUserByID(userClaims.ID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	measurement := &database.Measurement{
		ID:     in.ID,
		UserID: User.Id,
		Date:   in.Date,
		Value:  in.Value,
	}
	updatedMeasurement, err := database.UpdateMeasurement(measurement)
	if err != nil {
		return nil, err
	}
	return &UpdateMeasurementOutput{
		Measurement: *updatedMeasurement,
	}, nil
}

type DeleteMeasurementInput struct {
	ID int `json:"id" path:"id" binding:"required"`
}

type DeleteMeasurementOutput struct {
	Status string `json:"status"`
}

func DeleteMeasurement(c *gin.Context, in *DeleteMeasurementInput) (*DeleteMeasurementOutput, error) {
	if err := database.DeleteMeasurement(in.ID); err != nil {
		return nil, err
	}
	return &DeleteMeasurementOutput{
		Status: "measurement deleted successfully",
	}, nil
}

type AddTrainInput struct {
	// UserID    int    `json:"userId" binding:"required"`
	Date      string `json:"date" binding:"required"`
	TrainerID int    `json:"trainerId" binding:"required"`
	ClientID  int    `json:"clientId" binding:"required"`
	Duration  string `json:"duration" binding:"required"`
}

type AddTrainOutput struct {
	Train database.Train `json:"train"`
}

func AddTrain(c *gin.Context, in *AddTrainInput) (*AddTrainOutput, error) {
	claims := c.MustGet("claims").(jwt.MapClaims)

	userClaims, err := auth.ExtractClaims(claims)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	User, err := database.FindUserByID(userClaims.ID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	train := &database.Train{
		UserID:    User.Id,
		Date:      in.Date,
		TrainerID: in.TrainerID,
		ClientID:  in.ClientID,
		Duration:  in.Duration,
	}
	createdTrain, err := database.AddTrain(train)
	if err != nil {
		return nil, err
	}
	return &AddTrainOutput{
		Train: *createdTrain,
	}, nil
}

type UpdateTrainInput struct {
	ID int `json:"id" path:"id" binding:"required"`
	// UserID    int    `json:"userId" binding:"required"`
	Date      string `json:"date" binding:"required"`
	TrainerID int    `json:"trainerId" binding:"required"`
	ClientID  int    `json:"clientId" binding:"required"`
	Duration  string `json:"duration" binding:"required"`
}

type UpdateTrainOutput struct {
	Train database.Train `json:"train"`
}

func UpdateTrain(c *gin.Context, in *UpdateTrainInput) (*UpdateTrainOutput, error) {
	claims := c.MustGet("claims").(jwt.MapClaims)

	userClaims, err := auth.ExtractClaims(claims)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	User, err := database.FindUserByID(userClaims.ID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	train := &database.Train{
		ID:        in.ID,
		UserID:    User.Id,
		Date:      in.Date,
		TrainerID: in.TrainerID,
		ClientID:  in.ClientID,
		Duration:  in.Duration,
	}
	updatedTrain, err := database.UpdateTrain(train)
	if err != nil {
		return nil, err
	}
	return &UpdateTrainOutput{
		Train: *updatedTrain,
	}, nil
}

type DeleteTrainInput struct {
	ID int `json:"id" path:"id" binding:"required"`
}

type DeleteTrainOutput struct {
	Status string `json:"status"`
}

func DeleteTrain(c *gin.Context, in *DeleteTrainInput) (*DeleteTrainOutput, error) {
	if err := database.DeleteTrain(in.ID); err != nil {
		return nil, err
	}
	return &DeleteTrainOutput{
		Status: "train deleted successfully",
	}, nil
}
