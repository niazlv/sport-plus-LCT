package user

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/juju/errors"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/niazlv/sport-plus-LCT/internal/api/auth"
	database "github.com/niazlv/sport-plus-LCT/internal/database/auth"
	"github.com/wI2L/fizz"
)

type JWTAuthSecurityType struct {
	Name         string
	Type         string
	BearerFormat string
	Scheme       string
}

// example authenticatedHandler(v1, "GET", "/protected", "Protected endpoint", getProtected)
// func authenticatedHandler(group *fizz.RouterGroup, method, path, summary string, handler interface{}) {
// 	group.Handle(method, path, []fizz.OperationOption{
// 		fizz.Summary(summary),
// 		fizz.Security(&openapi.SecurityRequirement{
// 			"BearerAuth": {},
// 		}),
// 	}, tonic.Handler(handler, 200))
// }

func Setup(rg *fizz.RouterGroup) {
	api := rg.Group("user", "User", "User related endpoints")

	_ = api
	api.GET("", []fizz.OperationOption{fizz.Summary("Return User"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetUser, 200))

	//api.GET("", []fizz.OperationOption{fizz.Summary("Check auth status")}, WithAuth, tonic.Handler(getGet, 200))
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
