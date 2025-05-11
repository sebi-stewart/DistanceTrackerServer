package constants

import "os"

var (
	ServerPort   = os.Getenv("DTS_PORT")
	DatabaseFile = os.Getenv("DTS_DB_FILE")
)
