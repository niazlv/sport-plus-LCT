package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/niazlv/sport-plus-LCT/internal/config"
	database "github.com/niazlv/sport-plus-LCT/internal/database/auth"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var db *gorm.DB

var secretKey = []byte("my-secret-key-public")

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Setup(rg *gin.RouterGroup) {
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

	// defer db.Close()

	// Users = []database.User{
	// 	{Id: 0, Login: "test", Password: "123"},
	// 	{Id: 1, Login: "test@g.co", Password: "123"},
	// }

	api := rg.Group("auth")
	api.GET("", WithAuth, getGet)
	api.GET("/signin", getSignin)
	api.POST("/signup", postSignup)
}

func getGet(c *gin.Context) {
	c.String(http.StatusOK, "test, claims:\n", c.MustGet("claims").(jwt.MapClaims))
}

func getSignin(c *gin.Context) {
	passwd := c.Request.URL.Query().Get("password")
	login := c.Request.URL.Query().Get("login")
	if passwd == "" || login == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "login or password can't be null"})
		return
	}
	log.Println("login: ", login)
	log.Println("passwd: ", passwd)

	User, err := database.FindUserByLogin(login)
	if err != nil {
		log.Println("ERROR getSignin(): ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DATABASE ERROR"})
		return
	}

	if User.Login == login && (CheckPasswordHash(passwd, User.Password) || passwd == User.Password) {
		token := createToken(User)
		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user":  User,
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "user with this login and password not found!"})
}

func postSignup(c *gin.Context) {
	var creds database.User
	err := json.NewDecoder(c.Request.Body).Decode(&creds)
	if err != nil || (creds == database.User{}) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "credencials can't be null"})
		return
	}
	if creds.Login == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "login can't be null"})
		return
	}
	if creds.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password can't be null"})
		return
	}

	User, err := database.FindUserByLogin(creds.Login)
	if err != nil {
		log.Println("ERROR postSignup(): ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DATABASE ERROR"})
		return
	}
	if User != nil && User.Login == creds.Login {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user with this login is already created"})
		return
	}
	creds.Password, err = HashPassword(creds.Password)

	if err != nil {
		log.Println("ERROR postSignup(), hashing password: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing error!"})
		return
	}

	user := database.User{
		Login:    creds.Login,
		Password: creds.Password,
		Role:     creds.Role, //1 - trainer ,0 - User
	}

	_, err = database.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	token := createToken(&user)
	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
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
