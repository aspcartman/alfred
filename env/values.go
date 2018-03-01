package env

import "os"

func TelegramToken() string {
	return os.Getenv("TELEGRAM_TOKEN")
}