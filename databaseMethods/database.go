package databaseMethods

import (
	"DiscordBot/pkg/logger/logger"
	"github.com/bwmarrin/discordgo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

// Nicknames - таблица изменений никнеймов
type Nicknames struct {
	gorm.Model
	Nickname   string    `gorm:"type:varchar(255);not null"`
	DateChange time.Time `gorm:"type:datetime;not null"`
}

// BotUsage - таблица использования бота
type BotUsage struct {
	gorm.Model
	UserGlobalName string    `gorm:"type:varchar(255);not null"`
	Command        string    `gorm:"type:varchar(255);not null"`
	DateUsage      time.Time `gorm:"type:datetime;not null"`
}

// OpenDatabase - подключение к базе данных
func OpenDatabase(dbPath string, log *logger.Log) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Error("ошибка базы данных", logger.GetPlace())
		return nil, err
	}
	if !db.Migrator().HasTable(&Nicknames{}) {
		err = db.AutoMigrate(&Nicknames{})
		if err != nil {
			log.Error("Ошибка создания таблицы SergeyNicknames", logger.GetPlace())
			return nil, err
		}
	}
	if !db.Migrator().HasTable(&BotUsage{}) {
		err = db.AutoMigrate(&BotUsage{})
		if err != nil {
			log.Error("Ошибка создания таблицы BotUsage", logger.GetPlace())
			return nil, err
		}
	}
	return db, nil
}

// ChangeNickname - Смена ника Сереге
func ChangeNickname(s *discordgo.Session, GuildId, ChannelId, UserId, newNickname string, db *gorm.DB, logs *logger.Log) {
	db.Create(&Nicknames{
		Nickname:   newNickname,
		DateChange: time.Now(),
	})
	s.GuildMemberNickname(GuildId, UserId, newNickname)
	logs.Info("Ник Сереги успешно изменен на "+newNickname, logger.GetPlace())
}

func DBNewAction(User, Message string, db *gorm.DB, logs *logger.Log) {
	db.Create(&BotUsage{
		UserGlobalName: User,
		Command:        Message,
		DateUsage:      time.Now(),
	})
	logs.Info(User+" использовал команду <"+Message+">", logger.GetPlace())
}
