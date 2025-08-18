package databaseMethods

import (
	"DiscordBot/pkg/logger/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type SergeyNickname struct {
	gorm.Model
	Nickname   string    `gorm:"type:varchar(255);not null"`
	DateChange time.Time `gorm:"type:datetime;not null"`
}

type BotUsage struct {
	gorm.Model
	UserGlobalName string    `gorm:"type:varchar(255);not null"`
	Command        string    `gorm:"type:varchar(255);not null"`
	DateUsage      time.Time `gorm:"type:datetime;not null"`
}

func OpenDatabase(dbPath string, log *logger.Log) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Error("ошибка базы данных", logger.GetPlace())
		return nil, err
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()
	if !db.Migrator().HasTable(&SergeyNickname{}) {
		err = db.AutoMigrate(&SergeyNickname{})
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
