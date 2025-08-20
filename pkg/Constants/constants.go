package Constants

import "path/filepath"

var (
	PathToNicknamestxt = filepath.Join("..", "..", "AI", "Nicknames.txt")
	PathToBotSystemtxt = filepath.Join("..", "..", "AI", "BotsystemPromt.txt")
	PathToDataBasetxt  = filepath.Join("..", "..", "databaseMethods", "database", "database.db")
)

const TalksOnlyInServer = "Йоу, я отвечаю только на сервере 'не придумал', "
