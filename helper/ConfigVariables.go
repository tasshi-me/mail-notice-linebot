package helper

import "os"

// ConfigVars ..
func ConfigVars() ConfigVariables {
	return ConfigVariables{
		EnvLoaded:    os.Getenv("ENV_LOADED"),
		DocumentRoot: os.Getenv("DOCUMENT_ROOT"),

		Port:   os.Getenv("PORT"),
		DBPort: os.Getenv("DBPORT"),

		MongodbURI: os.Getenv("MONGODB_URI"),

		HerokuAppName: os.Getenv("HEROKU_APP_NAME"),

		LineAPI: LineAPIConfigVariables{
			ChannelID:     os.Getenv("LINE_CHANNEL_ID"),
			ChannelSecret: os.Getenv("LINE_CHANNEL_SECRET"),
			AccessToken:   os.Getenv("LINE_ACCESS_TOKEN"),
		},

		SMTP: SMTPConfigVariables{
			SenderUsername: os.Getenv("SENDER_USERNAME"),
			SenderAddress:  os.Getenv("SENDER_ADDRESS"),
			ServerName:     os.Getenv("SMTP_SERVER_NAME"),
			AuthUser:       os.Getenv("SMTP_AUTH_USER"),
			AuthPassword:   os.Getenv("SMTP_AUTH_PASSWORD"),
		},

		IMAP: IMAPConfigVariables{
			Address:      os.Getenv("IMAP_ADDRESS"),
			ServerName:   os.Getenv("IMAP_SERVER_NAME"),
			AuthUser:     os.Getenv("IMAP_AUTH_USER"),
			AuthPassword: os.Getenv("IMAP_AUTH_PASSWORD"),
			MboxName:     os.Getenv("IMAP_MBOX_NAME"),
		},
	}
}

// ConfigVariables ..
type ConfigVariables struct {
	EnvLoaded    string
	DocumentRoot string

	Port   string
	DBPort string

	MongodbURI string

	HerokuAppName string

	LineAPI LineAPIConfigVariables
	SMTP    SMTPConfigVariables
	IMAP    IMAPConfigVariables
}

// LineAPIConfigVariables ..
type LineAPIConfigVariables struct {
	ChannelID     string
	ChannelSecret string
	AccessToken   string
}

// SMTPConfigVariables ..
type SMTPConfigVariables struct {
	SenderUsername string
	SenderAddress  string
	ServerName     string
	AuthUser       string
	AuthPassword   string
}

// IMAPConfigVariables ..
type IMAPConfigVariables struct {
	Address      string
	ServerName   string
	AuthUser     string
	AuthPassword string
	MboxName     string
}
