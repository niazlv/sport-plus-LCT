// internal/database/chat/chat.go
package chat

import (
	"fmt"
	"time"

	"github.com/niazlv/sport-plus-LCT/internal/config"
	"github.com/niazlv/sport-plus-LCT/internal/database/auth"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AttachableType string

type Attachment struct {
  Id int `gorm:"primaryKey" json:"id"`
  AttachableType AttachableType `json:"attachable_type"`
  AttachableId int `json:"attachable_id"`
  MessageId int `json:"message_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	Id        int       `gorm:"primaryKey" json:"id"`
	ChatId    int       `json:"chat_id"`
	UserId    int       `json:"user_id"`
	Content   string    `json:"content"`
  Attachments []*Attachment `json:"attachments" gorm:"foreignKey:MessageId"`
	CreatedAt time.Time `json:"created_at"`
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

	err = db.AutoMigrate(&auth.Chat{}, &Message{}, &Attachment{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateChat(chat *auth.Chat) (*auth.Chat, error) {
	result := db.Create(chat)
	if result.Error != nil {
		return nil, result.Error
	}
	return chat, nil
}

func CanJoinChat(chatId int, userId int) bool {
    var count int64
    err := db.Model(&auth.Chat{}).
        Joins("JOIN chat_users ON chat_users.chat_id = chats.id").
        Where("chats.id = ? AND chat_users.user_id = ?", chatId, userId).
        Count(&count).Error

    if err != nil {
        return false
    }

    return count == 0
}

func GetChatByID(id int) (*auth.Chat, error) {
	var chat auth.Chat
	result := db.First(&chat, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &chat, nil
}

func GetChats() ([]auth.Chat, error) {
	var chats []auth.Chat
	result := db.Find(&chats)
	if result.Error != nil {
		return nil, result.Error
	}
	return chats, nil
}

type CreateMessageDto struct {
  Message Message
  AttachableId *int
  AttachableType *AttachableType
}

func CreateMessage(dto *CreateMessageDto) (*Message, error) {
    tx := db.Begin()
    if tx.Error != nil {
        return nil, tx.Error
    }

    if err := tx.Create(&dto.Message).Error; err != nil {
        tx.Rollback()
        return nil, err
    }

    attachment := &Attachment{
        AttachableType: *dto.AttachableType,
        AttachableId:   *dto.AttachableId,
        MessageId:      dto.Message.Id,
    }

    if err := tx.Create(attachment).Error; err != nil {
        tx.Rollback()
        return nil, err
    }

    dto.Message.Attachments = []*Attachment{attachment}

    if err := tx.Commit().Error; err != nil {
        return nil, err
    }

    return &dto.Message, nil
}

func GetMessagesByChatID(chatID int) ([]Message, error) {
	var messages []Message
	result := db.Preload("Attachments").Where("chat_id = ?", chatID).Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}
