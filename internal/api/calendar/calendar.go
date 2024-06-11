// internal/api/calendar/calendar.go
package calendar

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/niazlv/sport-plus-LCT/internal/api/auth"
	database "github.com/niazlv/sport-plus-LCT/internal/database/auth"
	"github.com/niazlv/sport-plus-LCT/internal/database/calendar"
	"github.com/wI2L/fizz"
)

// ScheduleInput структура для ввода данных для расписания
type ScheduleInput struct {
	ID             string    `path:"schedule_id" binding:"required"`
	ClientID       int       `json:"client_id"`
	Date           time.Time `json:"date"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	Type           string    `json:"type"`
	ReminderClient bool      `json:"reminder_client"`
	ReminderCoach  bool      `json:"reminder_coach"`
	IsGlobal       bool      `json:"is_global"`
}

func Setup(rg *fizz.RouterGroup) {
	api := rg.Group("calendar", "Calendar", "Calendar related endpoints")

	api.GET("", []fizz.OperationOption{fizz.Summary("Get schedules"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetSchedules, 200))
	api.GET("/global", []fizz.OperationOption{fizz.Summary("Get global schedules"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetGlobalSchedules, 200))
	api.GET("/local", []fizz.OperationOption{fizz.Summary("Get local schedules"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetLocalSchedules, 200))
	api.GET("/:schedule_id", []fizz.OperationOption{fizz.Summary("Get schedule by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(GetScheduleByID, 200))
	api.POST("", []fizz.OperationOption{fizz.Summary("Create a new schedule"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(CreateSchedule, 201))
	api.PUT("/:schedule_id", []fizz.OperationOption{fizz.Summary("Update schedule by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(UpdateSchedule, 200))
	api.DELETE("/:schedule_id", []fizz.OperationOption{fizz.Summary("Delete schedule by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(DeleteSchedule, 204))
}

func GetSchedules(c *gin.Context) (*[]calendar.Schedule, error) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	userClaims, err := auth.ExtractClaims(claims)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	User, err := database.FindUserByID(userClaims.ID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	var schedules []calendar.Schedule
	if User.Role == 1 { // Тренер
		schedules, err = calendar.GetSchedulesByCoachID(userClaims.ID)
	} else { // Пользователь
		schedules, err = calendar.GetSchedulesByClientID(userClaims.ID)
	}

	if err != nil {
		return nil, err
	}

	return &schedules, nil
}

func GetGlobalSchedules(c *gin.Context) (*[]calendar.Schedule, error) {
	schedules, err := calendar.GetGlobalSchedules()
	if err != nil {
		return nil, err
	}

	return &schedules, nil
}

func GetLocalSchedules(c *gin.Context) (*[]calendar.Schedule, error) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	userClaims, err := auth.ExtractClaims(claims)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	User, err := database.FindUserByID(userClaims.ID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	var schedules []calendar.Schedule
	if User.Role == 1 { // Тренер
		schedules, err = calendar.GetLocalSchedulesByCoachID(userClaims.ID)
	} else { // Пользователь
		schedules, err = calendar.GetLocalSchedulesByClientID(userClaims.ID)
	}

	if err != nil {
		return nil, err
	}

	return &schedules, nil
}

type GetScheduleByIDParams struct {
	ID string `path:"schedule_id" binding:"required"`
}

func GetScheduleByID(c *gin.Context, params *GetScheduleByIDParams) (*calendar.Schedule, error) {
	idStr := params.ID
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, errors.New("invalid schedule_id")
	}

	schedule, err := calendar.GetScheduleByID(id)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

func CreateSchedule(c *gin.Context, in *ScheduleInput) (*calendar.Schedule, error) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	userClaims, err := auth.ExtractClaims(claims)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	newSchedule := calendar.Schedule{
		CoachID:        userClaims.ID,
		ClientID:       in.ClientID,
		Date:           in.Date,
		StartTime:      in.StartTime,
		EndTime:        in.EndTime,
		Type:           in.Type,
		ReminderClient: in.ReminderClient,
		ReminderCoach:  in.ReminderCoach,
		IsGlobal:       in.IsGlobal,
	}

	savedSchedule, err := calendar.CreateSchedule(&newSchedule)
	if err != nil {
		return nil, err
	}

	return savedSchedule, nil
}

func UpdateSchedule(c *gin.Context, in *ScheduleInput) (*calendar.Schedule, error) {
	idStr := in.ID
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, errors.New("invalid schedule_id")
	}

	schedule, err := calendar.GetScheduleByID(id)
	if err != nil {
		return nil, err
	}

	if schedule == nil {
		return nil, errors.New("schedule not found")
	}

	schedule.ClientID = in.ClientID
	schedule.Date = in.Date
	schedule.StartTime = in.StartTime
	schedule.EndTime = in.EndTime
	schedule.Type = in.Type
	schedule.ReminderClient = in.ReminderClient
	schedule.ReminderCoach = in.ReminderCoach
	schedule.IsGlobal = in.IsGlobal

	err = calendar.UpdateSchedule(schedule)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

type DeleteScheduleParams struct {
	ID string `path:"schedule_id" binding:"required"`
}

func DeleteSchedule(c *gin.Context, params *DeleteScheduleParams) error {
	idStr := params.ID
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errors.New("invalid schedule_id")
	}

	err = calendar.DeleteSchedule(id)
	if err != nil {
		return err
	}

	return nil
}