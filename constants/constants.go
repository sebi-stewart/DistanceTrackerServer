package constants

import "os"

var (
	ServerPort    = os.Getenv("DTS_PORT")
	DatabaseFile  = os.Getenv("DTS_DB_FILE")
	JWT_SecretKey = os.Getenv("DTS_JWT_SECRET_KEY")
)
