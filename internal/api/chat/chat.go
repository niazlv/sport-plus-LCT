package chat

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/niazlv/sport-plus-LCT/internal/api/auth"
	database_auth "github.com/niazlv/sport-plus-LCT/internal/database/auth"
	database "github.com/niazlv/sport-plus-LCT/internal/database/chat"
	"github.com/wI2L/fizz"
)

func Setup(rg *fizz.RouterGroup) {
	api := rg.Group("chat", "Chat", "Chat related endpoints")

	api.POST("", []fizz.OperationOption{fizz.Summary("Create a new chat"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(createChat, 200))
	api.GET("", []fizz.OperationOption{fizz.Summary("Get all chats"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(getChats, 200))
	api.GET("/:id", []fizz.OperationOption{fizz.Summary("Get a chat by ID"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(getChatByID, 200))
	// api.POST("/:id/message", []fizz.OperationOption{fizz.Summary("Send a message in a chat"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(sendMessage, 200))
	api.GET("/:id/messages", []fizz.OperationOption{fizz.Summary("Get all messages in a chat"), auth.BearerAuth}, auth.WithAuth, tonic.Handler(getMessages, 200))
}

type createChatInput struct {
	Name string `json:"name" validate:"required"`
}

type createChatOutput struct {
	Chat database_auth.Chat `json:"chat"`
}

func createChat(c *gin.Context, in *createChatInput) (*createChatOutput, error) {
	chat := &database_auth.Chat{
		Name:      in.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdChat, err := database.CreateChat(chat)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &createChatOutput{Chat: *createdChat}, nil
}

type getChatsOutput struct {
	Chats []database_auth.Chat `json:"chats"`
}

func getChats(c *gin.Context) (*getChatsOutput, error) {
	chats, err := database.GetChats()
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &getChatsOutput{Chats: chats}, nil
}

type getChatByIDInput struct {
	ID int `json:"id" path:"id" validate:"required"`
}

type getChatByIDOutput struct {
	Chat database_auth.Chat `json:"chat"`
}

func getChatByID(c *gin.Context, in *getChatByIDInput) (*getChatByIDOutput, error) {
	chat, err := database.GetChatByID(in.ID)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &getChatByIDOutput{Chat: *chat}, nil
}

type getMessagesOutput struct {
	Messages []database.Message `json:"messages"`
}

type getMessagesParam struct {
	ChatID int `path:"id" validate:"required"`
}

func getMessages(c *gin.Context, params *getMessagesParam) (*getMessagesOutput, error) {

	messages, err := database.GetMessagesByChatID(params.ChatID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &getMessagesOutput{Messages: messages}, nil
}
