package settings

const (
	ServerHost = "localhost"
	ServerPort = ":10001"
	ServerUrl  = "http://" + ServerHost + ServerPort

	SecretKey = "rpi service jwt"

	StaticRoot = "/static"
	StaticDir  = "static/"

	SMTPHost     = "smtp.126.com"
	SMTPAddress  = SMTPHost + ":25"
	SMTPUsername = "XXXXXXX"
	SMTPPassword = "XXXXXXXX"
)
