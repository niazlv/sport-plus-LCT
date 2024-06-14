package routes

import (
	"github.com/niazlv/sport-plus-LCT/internal/api/auth"
	"github.com/niazlv/sport-plus-LCT/internal/api/calendar"
	"github.com/niazlv/sport-plus-LCT/internal/api/chat"
	"github.com/niazlv/sport-plus-LCT/internal/api/course"
	"github.com/niazlv/sport-plus-LCT/internal/api/upload"
	"github.com/niazlv/sport-plus-LCT/internal/api/user"
	"github.com/wI2L/fizz"
)

func Setup(f *fizz.Fizz) {
	// Создаем группу маршрутов
	api := f.Group("/v1", "API v1", "API version 1")

	// Настраиваем маршруты для авторизации и пользователей
	auth.Setup(api)
	user.Setup(api)
	course.Setup(api)
	calendar.Setup(api)
	upload.SetupUploadRoutes(api)
	chat.Setup(api)
}
