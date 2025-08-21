package cmd

import (
	"DiscordBot/databaseMethods"
	"DiscordBot/pkg/logger/logger"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"io"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

const MyServerId = "537698381527777300"
const SergeyId = "664192938460446730"
const LastNicknameChanges = 3
const TextChannelID = "904769583540486164"
const HoursToChangeNickname = 2
const HoursToCheckOfNicknameChanger = 1

var (
	lastExecution = time.Now() // для смены ника
	mu            sync.Mutex   // для бд
)

// NicknamesChanger - Функция, которая раз в 2 дня меняет ник сереге
func NicknamesChanger(s *discordgo.Session, UserId string, Nicknames []string, db *gorm.DB, logs *logger.Log) {
	for {
		if UserId != SergeyId {
			return
		}
		// Смотрим прошло ли 2 дня с момента смены ника
		if time.Since(lastExecution).Hours() >= HoursToChangeNickname {
			lastExecution = time.Now()
			// берем из базы данных три последних изменения ников
			var nicks []databaseMethods.Nicknames
			db.Order("id desc").Limit(LastNicknameChanges).Find(&nicks)
			if len(nicks) > 0 {
				// если бд не пустая
				newNickname := "Серега"
				for {
					// смотрим, чтобы ник бьл не повторяющимся
					nick := Nicknames[rand.Intn(len(Nicknames))]
					ok := 0
					for _, v := range nicks {
						if v.Nickname != nick {
							ok++
						}
					}
					if ok == len(nicks) {
						newNickname = nick
						break
					}
				}
				// меняем
				databaseMethods.ChangeNickname(s, MyServerId, TextChannelID, UserId, newNickname, db, logs)
			} else { // Если бд пустая
				newNick := Nicknames[rand.Intn(len(Nicknames))]
				databaseMethods.ChangeNickname(s, MyServerId, TextChannelID, UserId, newNick, db, logs)
			}
		}
		// делаем перерыв
		time.Sleep(HoursToCheckOfNicknameChanger)
	}
}

// GetNicknames - получает из txt файла никнеймы сереги
func GetNicknames(path string, logs *logger.Log) ([]string, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logs.Error("Nicknames File is not exist", logger.GetPlace())
		return nil, err
	}
	defer file.Close()
	bytes, _ := io.ReadAll(file)
	str := string(bytes)
	return strings.Split(str, "\n"), nil
}
