package main

import (
	"github.com/aspcartman/alfred/telegram"
	"github.com/aspcartman/alfred/env"
	"time"
	"github.com/aspcartman/alfred/torrent"
	"github.com/aspcartman/exceptions"
)

func main() {
	env.Log.Info("starting")

	tor := torrent.NewBot("127.0.0.1:8083")
	env.Log.Info("torrent bot started")

	go runbot(tor)



	time.Sleep(1 * time.Hour)
}

func runbot(tor *torrent.Bot) {
	defer e.Catch(func(e *e.Exception) {
		go runbot(tor)
	})

	telegram.RunBot(env.TelegramToken(), func(s *telegram.Session) {
		env.Log.Info("msg in")
		s.Reply("ща ща...")
		msg := s.GetMessage()
		switch {
		case msg.Document != nil:
			tor.Handle(s, msg)
		default:
			s.Reply("не осилил =(")
		}
	})
	env.Log.Info("telegram started")

}