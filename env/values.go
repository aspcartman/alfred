package env

import (
	"os"
	"github.com/aspcartman/exceptions"
	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("Environment variable is empty")

func TelegramToken() string {
	return getkey("TELEGRAM_TOKEN")
}


func getkey(key string) string {
	token := os.Getenv(key)
	if len(token) == 0 {
		e.Throw("Key not found", ErrNotFound, e.Map{
			"key": key,
		})
	}
	return token
}