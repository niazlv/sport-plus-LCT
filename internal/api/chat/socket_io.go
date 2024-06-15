package chat

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/niazlv/sport-plus-LCT/internal/api/auth"
	database_auth "github.com/niazlv/sport-plus-LCT/internal/database/auth"
	database "github.com/niazlv/sport-plus-LCT/internal/database/chat"
)

var Server *socketio.Server

func GetRoomName(chatId int) string {
	return fmt.Sprintf(`room-%s`, strconv.Itoa(chatId))
}

func InitSocketIO() {
	Server = socketio.NewServer(nil)

	Server.OnConnect("/", func(s socketio.Conn) error {
		log.Println("connected:", s.ID())
		return nil
	})

	Server.OnEvent("/", "join_chat", func(s socketio.Conn, chatID int) {
		// Извлекаем токен из заголовка
		tokenStr := s.RemoteHeader().Get("Authorization")
		if tokenStr == "" {
			s.Close()
			return
		}
		// Удаляем префикс "Bearer "
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		// Проверяем валидность токена
		valid, err := auth.IsValidToken(tokenStr)
		if !valid || err != nil {
			s.Close()
			return
		}

		// Извлекаем claims из токена
		claims, err := auth.ExtractClaimsFromJWT(tokenStr)
		if err != nil {
			s.Close()
			return
		}

		userID, ok := claims["id"].(float64)
		if !ok {
			s.Close()
			return
		}

		// Находим пользователя по ID
		User, err := database_auth.FindUserByID(int(userID))
		if err != nil {
			s.Close()
			return
		}

		canJoin := database.CanJoinChat(chatID, User.Id)
		if !canJoin {
			return
		}
		s.Join(GetRoomName(chatID))

		log.Printf("user %s joined chat %d", s.ID(), chatID)
	})

	Server.OnEvent("/", "message", func(s socketio.Conn, dto database.CreateMessageDto) {
		dto.Message.CreatedAt = time.Now()
		createdMessage, err := database.CreateMessage(&dto)
		if err != nil {
			s.Emit("error", err.Error())
			return
		}

		// Broadcast the message to all users in the chat
		Server.BroadcastToRoom("/", GetRoomName(dto.Message.ChatId), "message", createdMessage)
	})

	Server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	Server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})

	go Server.Serve()
}

func SocketIOHandler(c *gin.Context) {
	Server.ServeHTTP(c.Writer, c.Request)
}
