package user

import (
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

	return &putOnboardingOutput{Status: "user updated successfully"}, nil
}
