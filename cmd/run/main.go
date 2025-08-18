package main

import (
	"DiscordBot/AI"
	"DiscordBot/cmd"
	"DiscordBot/pkg/Error"
	"DiscordBot/pkg/logger/logger"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const MyServerId = "537698381527777300"

// TODO: реализовать базу данных(использование бота, ники сереги), может туда еще логи пихнуть
func main() {
	logs := logger.NewLog()
	// Настраиваем переменные среды
	AIApi := os.Getenv("AI_API_KEY")
	if AIApi == "" {
		logs.Error(Error.ApiKeyIsEmpty, logger.GetPlace())
		panic("AI_API_KEY environment variable not set")
	}
	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	if botToken == "" {
		logs.Error(Error.BotTokenIsEmpty, logger.GetPlace())
		panic("DISCORD_BOT_TOKEN environment variable not set")
	}

	// Создаем сессию Discord
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		logs.Error(Error.SessionError+"\n"+err.Error(), logger.GetPlace())
		return
	}

	// Достаем системный промт
	file, errFile := os.ReadFile("../../AI/BotsystemPromt.txt")
	if errFile != nil {
		logs.Error(Error.SystemPromtFileDoesNotOpen+"\n"+errFile.Error(), logger.GetPlace())
		panic(errFile)
	}
	systemPromt := string(file)

	// Обработчик события "готовности" бота
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logs.Info(fmt.Sprintf("Бот запущен как %s#%s", r.User.Username, r.User.Discriminator), logger.GetPlace())
	})

	// Slash-команды
	commands := cmd.GetBotsCommands()

	// Обработчик Slash-команд
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		cmds := i.ApplicationCommandData()
		switch cmds.Name {
		case "ping":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Я обитаю независимо от твоего понимания. Ну вроде пишу тебе хуйню какую-то.",
				},
			})
		case "you":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Я - великий дух сервера 'не придумал', я меняю ник сереге, потому что на него сошла моя кара.",
				},
			})
		case "kara":
			// Указываем конкретный ID пользователя
			nicknames := []string{"Лошок", "Опездол", "Чевапчич", "Фидир", "Уебище", "Чушпан",
				"Пипипупу", "Черкашок", "Тупик", "3070м", "Глупи"}
			lenNicknames := len(nicknames)
			userID := "664192938460446730"
			newNick := "Серега " + nicknames[rand.Intn(lenNicknames-1)]

			err = s.GuildMemberNickname(i.GuildID, userID, newNick)
			if err != nil {
				logs.Warning("Не меняется ник у сереги", logger.GetPlace())
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Нет вайбика менять ник Сереге пока что",
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Ник Сереги был изменен на `%s`", newNick),
					},
				})
				logs.Info(fmt.Sprintf("Ник Сереги был изменен на `%s`", newNick), logger.GetPlace())
			}
		case "time":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Время по великому духу: " + time.Now().Format("02.01.2006 15:04:05"),
				},
			})
		case "talk":

			// 1. Немедленный отложенный ответ
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})
			if err != nil {
				logs.Warning(err.Error(), logger.GetPlace())
				return
			}

			// 2. Асинхронно обрабатываем запрос
			go func() {
				message := cmds.Options[0].StringValue()
				aiResponse, err := AI.Promt(message, systemPromt, AIApi)
				if err != nil {
					logs.Warning(err.Error(), logger.GetPlace())
					aiResponse = "Я чет устал пиздеть, идите нахуй"
				}

				// 3. Отправляем результат
				_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: aiResponse,
				})
				if err != nil {
					logs.Warning(err.Error(), logger.GetPlace())
				}
			}()
		default:
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Але, фигню мне тут не пиши. А то и на тебя моя кара падет",
				},
			})
		}
	})

	// Открываем соединение
	err = dg.Open()
	if err != nil {
		logs.Error(Error.SessionError+"\n"+err.Error(), logger.GetPlace())
		return
	}
	defer dg.Close() // Закрываем соединение при выходе

	// Регистрация команд
	registeredCommands, err := dg.ApplicationCommandBulkOverwrite(dg.State.User.ID, MyServerId, commands)
	if err != nil {
		logs.Warning(Error.RegisteringCommandsError+": "+err.Error(), logger.GetPlace())
		panic(Error.RegisteringCommandsError + ": " + err.Error())
	}
	log.Println("Зарегистрированные команды:", registeredCommands)

	// Ждем сигнала завершения (Ctrl+C)
	fmt.Println("Бот работает. Нажмите Ctrl+C для выхода.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
