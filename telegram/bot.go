package telegram

import (
	"time"

	"github.com/aspcartman/exceptions"
	"gopkg.in/tucnak/telebot.v2"
	"github.com/aspcartman/alfred/env"
	"net/http"
	"io/ioutil"
	"sync"
)

type HandlerFunc func(s *Session)

type Bot struct {
	*telebot.Bot
	sessions map[int]*Session
	handler  HandlerFunc
	lck      sync.Mutex
}

func RunBot(token string, handler HandlerFunc) *Bot {
	telegramBot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		Reporter: func(err error) {
			env.Log.WithError(err).Error("exception in tg bot handler")
		},
	})
	if err != nil {
		e.Throw("Failed to authorize Bot a instance", err, e.Map{
			"token": token,
		})
	}

	b := &Bot{telegramBot, map[int]*Session{}, handler, sync.Mutex{}}

	for _, event := range []string{telebot.OnText, telebot.OnDocument} {
		telegramBot.Handle(event, b.Handle)
	}

	go telegramBot.Start()

	return b
}

func (b *Bot) Handle(msg *telebot.Message) {
	b.lck.Lock()
	session, ok := b.sessions[msg.Sender.ID]
	if !ok {
		session = &Session{b, msg.Sender.ID, make(chan *telebot.Message, 10)}
		b.sessions[msg.Sender.ID] = session

		go func() {
			defer e.Catch(func(e *e.Exception) {
				// do nothing
			})

			defer func() {
				b.lck.Lock()
				delete(b.sessions, msg.Sender.ID)
				b.lck.Unlock()
			}()

			defer
			b.handler(session)

		}()

	}
	b.lck.Unlock()

	session.recv(msg)
}

func (b *Bot) Download(f string) []byte {
	url, err := b.FileURLByID(f)
	if err != nil {
		e.Throw("couldn't get a file url", err, e.Map{
			"file": f,
		})
	}

	res, err := http.Get(url)
	if err != nil {
		e.Throw("failed downloading file", err, e.Map{
			"url": url,
		})
	}
	if res.StatusCode != http.StatusOK {
		e.Throw("failed downloading file, bad status", nil, e.Map{
			"url":    url,
			"status": res.StatusCode,
		})
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		e.Throw("failed downloading file from body", err, e.Map{
			"url": url,
		})
	}

	return data
}
