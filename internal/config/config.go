package config

import (
	"os"
)

// Config структура для хранения конфигурации базы данных
type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	JWTSecret  string
}

// LoadConfig загружает конфигурацию из файла .env
func LoadConfig() (*Config, error) {
	// Загружаем основной файл .env
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Println("Error loading .env file")
	// 	return nil, err
	// }

	// // Определяем среду из переменной окружения
	// env := os.Getenv("APP_ENV")
	// if env == "" {
	// 	env = "dev"
	// }

	// // Загружаем файл соответствующей среды
	// envFile := fmt.Sprintf("env/%s.env", env)
	// err = godotenv.Load(envFile)
	// if err != nil {
	// 	log.Printf("Error loading %s file", envFile)
	// 	return nil, err
	// }

	config := &Config{
		DBUser:     os.Getenv("POSTGRES_USER"),
		DBPassword: os.Getenv("POSTGRES_PASSWORD"),
		DBName:     os.Getenv("POSTGRES_DB"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}

	return config, nil
}
