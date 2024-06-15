// internal/database/chat/chat.go
package chat

import (
	"fmt"
	"time"

	"github.com/niazlv/sport-plus-LCT/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Chat struct {
	Id        int       `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AttachableType string

type Attachment struct {
	Id             int            `gorm:"primaryKey" json:"id"`
	AttachableType AttachableType `json:"attachable_type"`
	AttachableId   int            `json:"attachable_id"`
	CreatedAt      time.Time      `json:"created_at"`
}

type Message struct {
	Id          int          `gorm:"primaryKey" json:"id"`
	ChatId      int          `json:"chat_id"`
	UserId      int          `json:"user_id"`
	Content     string       `json:"content"`
	Attachments []Attachment `json:"attachments"`
	CreatedAt   time.Time    `json:"created_at"`
}

var db *gorm.DB

func InitDB() (*gorm.DB, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort)

	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	if db == nil {
		return nil, fmt.Errorf("failed to connect to database")
	}

	err = db.AutoMigrate(&Chat{}, &Message{}, &Attachment{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateChat(chat *Chat) (*Chat, error) {
	result := db.Create(chat)
	if result.Error != nil {
		return nil, result.Error
	}
	return chat, nil
}

func GetChatByID(id int) (*Chat, error) {
	var chat Chat
	result := db.First(&chat, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &chat, nil
}

func GetChats() ([]Chat, error) {
	var chats []Chat
	result := db.Find(&chats)
	if result.Error != nil {
		return nil, result.Error
	}
	return chats, nil
}

func CreateMessage(message *Message) (*Message, error) {
	result := db.Create(message)
	if result.Error != nil {
		return nil, result.Error
	}
	return message, nil
}

func GetMessagesByChatID(chatID int) ([]Message, error) {
	var messages []Message
	result := db.Preload("Attachments").Where("chat_id = ?", chatID).Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}
