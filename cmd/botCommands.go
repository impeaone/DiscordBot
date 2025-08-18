package cmd

import "github.com/bwmarrin/discordgo"

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
