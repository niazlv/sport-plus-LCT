package chat

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	database "github.com/niazlv/sport-plus-LCT/internal/database/chat"
)

var Server *socketio.Server

func InitSocketIO() {
	Server = socketio.NewServer(nil)

	Server.OnConnect("/", func(s socketio.Conn) error {
		log.Println("connected:", s.ID())
		return nil
	})

	Server.OnEvent("/", "join_chat", func(s socketio.Conn, chatID int) {
		s.Join(string(rune(chatID)))
		log.Printf("user %s joined chat %d", s.ID(), chatID)
	})

	Server.OnEvent("/", "message", func(s socketio.Conn, msg database.Message) {
		msg.CreatedAt = time.Now()
		createdMessage, err := database.CreateMessage(&msg)
		if err != nil {
			s.Emit("error", err.Error())
			return
		}

		// Broadcast the message to all users in the chat
		Server.BroadcastToRoom("/", string(rune(msg.ChatId)), "message", createdMessage)
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
