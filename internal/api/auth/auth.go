package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/juju/errors"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/niazlv/sport-plus-LCT/internal/config"
	database "github.com/niazlv/sport-plus-LCT/internal/database/auth"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var db *gorm.DB

var secretKey = []byte("my-secret-key-public")

type UserClaims struct {
	Exp   int64  `json:"exp"`
	ID    int    `json:"id"`
	Login string `json:"login"`
}

var BearerAuth = fizz.Security(&openapi.SecurityRequirement{
	"BearerAuth": {},
})

func ExtractClaims(claims jwt.MapClaims) (*UserClaims, error) {
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid type for exp")
	}

	id, ok := claims["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid type for id")
	}

	login, ok := claims["login"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid type for login")
	}

	userClaims := &UserClaims{
		Exp:   int64(exp),
		ID:    int(id),
		Login: login,
	}

	return userClaims, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Setup(rg *fizz.RouterGroup) {
	var err error
	db, err = database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Load config error!", err)
	}

	// reload secretKey from config, if he exist's
	if cfg.JWTSecret != "" {
		secretKey = []byte(cfg.JWTSecret)
	}

	// Create a sub-group for auth routes
	api := rg.Group("/auth", "Auth", "Authentication related endpoints")

	// Define routes
	api.GET("", []fizz.OperationOption{fizz.Summary("Check auth status"), BearerAuth}, WithAuth, tonic.Handler(getGet, 200))
	api.GET("/signin", []fizz.OperationOption{fizz.Summary("Sign in")}, tonic.Handler(getSignin, 200))
	api.POST("/signup", []fizz.OperationOption{fizz.Summary("Sign up")}, tonic.Handler(postSignup, 200))
}

type getGetOutput struct {
	Message string `json:"message"`
}

func getGet(c *gin.Context) (*getGetOutput, error) {
	return &getGetOutput{
		Message: fmt.Sprint("test, claims:\n", c.MustGet("claims").(jwt.MapClaims)),
	}, nil
}

type getSigninOutput struct {
	Token string         `json:"token"`
	User  *database.User `json:"user"`
}

type getSigninInput struct {
	Login    string `json:"login" query:"login"`
	Password string `json:"password" query:"password"`
}

func getSignin(c *gin.Context, in *getSigninInput) (*getSigninOutput, error) {
	passwd := in.Password
	login := in.Login
	if passwd == "" || login == "" {
		return nil, errors.BadRequestf("login or password can't be null")
	}
	log.Println("login: ", login)
	log.Println("passwd: ", passwd)

	User, err := database.FindUserByLogin(login)
	if err != nil {
		log.Println("ERROR getSignin(): ", err)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "DATABASE ERROR"})
		return nil, fmt.Errorf("DATABASE ERROR")
	}

	if User.Login == login && (CheckPasswordHash(passwd, User.Password) || passwd == User.Password) {
		token := createToken(User)
		return &getSigninOutput{
			Token: token,
			User:  User,
		}, nil
	}

	return nil, errors.Unauthorizedf("user with this login and password not found!")
}

type postSignupInput struct {
	Login    string `json:"login" body:"login"`
	Password string `json:"password" body:"password"`
	Role     int    `json:"role" body:"role" default:"0"`
}

func postSignup(c *gin.Context, in *postSignupInput) (*getSigninOutput, error) {
	if (*in == postSignupInput{}) {
		return nil, errors.BadRequestf("credencials can't be null")
	}
	if in.Login == "" {
		return nil, errors.BadRequestf("login can't be null")
	}
	if in.Password == "" {
		return nil, errors.BadRequestf("password can't be null")
	}

	User, err := database.FindUserByLogin(in.Login)
	if err != nil {
		log.Println("ERROR postSignup(): ", err)
		return nil, fmt.Errorf("DATABASE ERROR")
	}
	if User != nil && User.Login == in.Login {
		return nil, errors.BadRequestf("user with this login is already created")
	}
	PasswordHashed, err := HashPassword(in.Password)

	if err != nil {
		log.Println("ERROR postSignup(), hashing password: ", err)
		return nil, fmt.Errorf("Password hashing error!")
	}

	user := database.User{
		Login:    in.Login,
		Password: PasswordHashed,
		Role:     in.Role, //1 - trainer ,0 - User
	}

	_, err = database.CreateUser(&user)
	if err != nil {
		return nil, errors.BadRequestf(err.Error())
	}
	log.Println("User", user)
	token := createToken(&user)
	return &getSigninOutput{
		User:  &user,
		Token: token,
	}, nil
}

// Создание JWT-токена
func createToken(user *database.User) string {
	claims := jwt.MapClaims{
		"login": user.Login,
		"id":    user.Id,
		"exp":   time.Now().Add(time.Hour * 24 * 180).Unix(), // Токен действителен 180 дней(пол года)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(secretKey)

	return signedToken
}

// Middleware для проверки токена
func WithAuth(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorizated"})
		c.Abort()
		return
	}

	// Проверяем, что токен начинается с "Bearer "
	if !strings.HasPrefix(tokenString, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат токена"})
		c.Abort()
		return
	}

	// Извлекаем сам токен
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token is expired"})
		c.Abort()
		return
	}

	// достаем значения
	claims := token.Claims.(jwt.MapClaims)
	log.Println(claims)
	User, err := database.FindUserByID(int(claims["id"].(float64)))
	log.Println(User)
	if err != nil {
		log.Println("ERROR WithAuth(): ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DATABASE ERROR"})
	}
	if User != nil {
		c.Set("claims", claims)
		c.Next()
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "user with this token not found"})
	c.Abort()
}
