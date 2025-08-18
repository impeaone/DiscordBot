package cmd

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

// GetBotsCommands - функция, возвращающая команды используемые ботом
func GetBotsCommands() []*discordgo.ApplicationCommand {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Проверка работоспособности бота",
		},
		{
			Name:        "you",
			Description: "Рассказывает о себе",
		},
		{
			Name:        "kara",
			Description: "Накладывает кару на серого, у него меняется никнейм",
		},
		{
			Name:        "time",
			Description: "Пишет текущее время. А хули, может кому-то надо",
		},
		{
			Name:        "talk",
			Description: "Отправить сообщение",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "Текст сообщения",
					Required:    true,
				},
			},
		},
	}
	return commands
}

// MessageForBot определяет, адресованно ли сообщение боту
func MessageForBot(message string) bool {
	ToLowerMessage := strings.ToLower(strings.TrimSpace(message))
	forBot := []string{"дух", "духа", "духу", "духом", "духе"}
	for _, v := range forBot {
		if strings.Contains(ToLowerMessage, v) {
			return true
		}
	}
	return false
}
