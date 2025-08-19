package main

import (
	"DiscordBot/AI"
	"DiscordBot/cmd"
	"DiscordBot/databaseMethods"
	"DiscordBot/pkg/Constants"
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

func main() {
	//Логи
	logs := logger.NewLog()
	// Объявляем RateLimiter для общения с духом
	RateLimiter := cmd.NewSimpleRateLimiter("", time.Now()) // Для ИИ от спама, можно писать ИИ раз в 5 секунд
	// Получаем никнеймы из txt
	nicknames, _ := cmd.GetNicknames(Constants.PathToNicknamestxt, logs)
	// Достаем системный промт
	systemPromt, _ := AI.GetSystemPromt(Constants.PathToBotSystemtxt, logs)
	// Подключаемся к бд
	db, err := databaseMethods.OpenDatabase(Constants.PathToDataBasetxt, logs)
	if err != nil {
		logs.Error(err.Error(), logger.GetPlace())
		return
	}
	// закрываем соединение с бд
	defer func() {
		sqldb, _ := db.DB()
		sqldb.Close()
	}()
	_ = db
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
	// запуск функции, которая каждые два дня меняет никнеймы сереге
	go cmd.NicknamesChanger(dg, cmd.SergeyId, nicknames, db, logs)

	// Обработчик события "готовности" бота
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logs.Info(fmt.Sprintf("Бот запущен как %s#%s", r.User.Username, r.User.Discriminator), logger.GetPlace())
	})

	// Slash-команды
	commands := cmd.GetBotsCommands()

	// Сообщения из чата
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Игнорируем сообщения от ботов
		if m.Author.Bot {
			return
		}
		// Диалог с духом без слеш-команды
		if cmd.MessageForBot(m.Content) {
			AiMessage, _ := AI.Promt(m.Author.GlobalName, m.Content, systemPromt, AIApi, RateLimiter)
			_, err := s.ChannelMessageSend(m.ChannelID, AiMessage)
			if err != nil {
				logs.Warning(Error.ChannelMessageError+"\n"+err.Error(), logger.GetPlace())
			}
		}
	})

	// Обработчик Slash-команд
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		cmds := i.ApplicationCommandData()
		switch cmds.Name {
		case "ping":
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Я обитаю независимо от твоего понимания. Ну вроде пишу тебе хуйню какую-то.",
				},
			})
			databaseMethods.DBNewAction(i.Interaction.Member.User.Username, cmds.Name, db, logs) // заносим событие в базу данных
			if err != nil {
				logs.Warning(Error.ChannelMessageError+"\n"+err.Error(), logger.GetPlace())
			}

		case "you":
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Я - великий дух сервера 'не придумал', я меняю ник сереге, потому что на него сошла моя кара.",
				},
			})
			databaseMethods.DBNewAction(i.Interaction.Member.User.Username, cmds.Name, db, logs) // заносим событие в базу данных
			if err != nil {
				logs.Warning(Error.ChannelMessageError+"\n"+err.Error(), logger.GetPlace())
			}

		case "kara":
			userID := cmd.SergeyId
			newNick := nicknames[rand.Intn(len(nicknames))]
			// Меняем ник
			err = s.GuildMemberNickname(i.GuildID, userID, newNick)
			if err != nil {
				logs.Warning(Error.ChangeNicknameError+"\n"+err.Error(), logger.GetPlace())
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Нет вайбика менять ник Сереге пока что",
					},
				})
				if err != nil {
					logs.Warning(Error.ChannelMessageError+"\n"+err.Error(), logger.GetPlace())
				}
			} else {
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Ник Сереги был изменен на `%s`", newNick),
					},
				})
				databaseMethods.ChangeNickname(s, cmd.MyServerId, cmd.TextChannelID, userID, newNick, db, logs)
				databaseMethods.DBNewAction(i.Interaction.Member.User.Username, cmds.Name, db, logs) // заносим событие в базу данных
				if err != nil {
					logs.Warning(Error.ChannelMessageError+"\n"+err.Error(), logger.GetPlace())
					return
				}
				logs.Info(fmt.Sprintf("Ник Сереги был изменен на `%s`", newNick), logger.GetPlace())
			}

		case "time":
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Время по великому духу: " + time.Now().Format("02.01.2006 15:04:05"),
				},
			})
			databaseMethods.DBNewAction(i.Member.User.Username, cmds.Name, db, logs) // заносим событие в базу данных
			if err != nil {
				logs.Warning(Error.ChannelMessageError+"\n"+err.Error(), logger.GetPlace())
			}

		case "talk":

			// Немедленный отложенный ответ (В дискорде появляется сообщение, что бот думает, он ждет ответа)
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})
			if err != nil {
				logs.Warning(Error.ChannelMessageError+"\n"+err.Error(), logger.GetPlace())
				return
			}

			// Асинхронно обрабатываем запрос
			go func() {
				message := cmds.Options[0].StringValue()
				aiResponse, err := AI.Promt(i.Member.User.GlobalName, message, systemPromt, AIApi, RateLimiter)
				if err != nil {
					logs.Warning(err.Error(), logger.GetPlace())
					aiResponse = "Все, Великий дух не хочет общаться"
				}

				// 3. Отправляем результат
				_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: aiResponse,
				})
				// заносим событие в базу данных
				databaseMethods.DBNewAction(i.Member.User.Username, cmds.Name+" "+message, db, logs) // заносим событие в базу данных
				if err != nil {
					logs.Warning(err.Error(), logger.GetPlace())
					return
				}
			}()
		default:
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Але, фигню мне тут не пиши. А то и на тебя моя кара падет",
				},
			})
			if err != nil {
				logs.Warning(Error.ChannelMessageError+"\n"+err.Error(), logger.GetPlace())
			}
		}
	})

	// Открываем соединение
	err = dg.Open()
	if err != nil {
		logs.Error(Error.SessionError+"\n"+err.Error(), logger.GetPlace())
		return
	}
	defer dg.Close()

	// Регистрация команд
	registeredCommands, err := dg.ApplicationCommandBulkOverwrite(dg.State.User.ID, cmd.MyServerId, commands)
	if err != nil {
		logs.Warning(Error.RegisteringCommandsError+": "+err.Error(), logger.GetPlace())
		panic(Error.RegisteringCommandsError + ": " + err.Error())
	}
	log.Println("Зарегистрированные команды:", registeredCommands)
	// Ждем сигнала завершения (Ctrl+C)
	fmt.Println("Бот работает. Ctrl+C для выхода.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
