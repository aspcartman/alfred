package telegram

import (
	"gopkg.in/tucnak/telebot.v2"
	"github.com/aspcartman/exceptions"
	"time"
	"errors"
	"strconv"
	"github.com/aspcartman/alfred/env"
)

var ErrSessionClosed = errors.New("session closed")

type Session struct {
	*Bot
	id    int
	queue chan *telebot.Message
}

func (s Session) Recipient() string {
	return strconv.Itoa(s.id)
}

func (s Session) recv(msg *telebot.Message) {
	s.queue <- msg
}

func (s Session) GetMessage() *telebot.Message {
	select {
	case msg := <-s.queue:
		return msg
	case <-time.After(1 * time.Minute):
	}
	e.Throw("session closed", ErrSessionClosed)
	return nil
}

func (s Session) Ask(question string, vars ...string) int {
	var btns [][]telebot.ReplyButton
	for _, v := range vars {
		btns = append(btns, []telebot.ReplyButton{{
			Text: v,
		}})
	}
	env.Log.WithField("vars", vars).Info("asking")

	s.Reply(question, &telebot.ReplyMarkup{
		ReplyKeyboard:   btns,
		OneTimeKeyboard: true,
	})

	msg := s.GetMessage()
	for i, s := range vars {
		if s == msg.Text {
			return i
		}
	}

	e.Throw("no such variant", nil, e.Map{
		"vars": vars,
		"answ": msg.Text,
	})

	return 0 // unreachable
}

func (s Session) Reply(obj interface{}, opts ...interface{}) *Response {
	msg, err := s.Send(s, obj, opts...)
	if err != nil {
		e.Throw("error replying", err)
	}
	return &Response{s, msg}
}
