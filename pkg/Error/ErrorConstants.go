package Error

const (
	LogFileDoesNotOpen         = "Log-file does not open"
	LogFileDoesNotWrite        = "Log-file does not write logs"
	SystemPromtFileDoesNotOpen = "System-promt-file does not open"
	NicknamesFileDoesNotOpen   = "Nicknames-file does not open"

	SessionError = "Create session error"

	ApiKeyIsEmpty   = "Api-key value in evironment is empty"
	BotTokenIsEmpty = "Bot-token value in evironment is empty"

	RegisteringCommandsError = "Registering commands error"

	ChannelMessageError = "Create channel message error"

	ChangeNicknameError = "Change nickname error"

	AiMessageError = "Ai message error"

	SessionLimit = "Session limit"
)

//TODO: тут много логов таких сделать надо
