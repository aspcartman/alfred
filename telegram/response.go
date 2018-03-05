package telegram

import (
	"gopkg.in/tucnak/telebot.v2"
	"github.com/aspcartman/exceptions"
	"strings"
)

type Response struct {
	session Session
	msg     *telebot.Message
}

func (r *Response) Edit(obj interface{}) {
	msg, err := r.session.Edit(r.msg, obj)
	if err != nil && !strings.Contains(err.Error(), "message is not modified") {
		e.Throw("failed editing response", err)
	} else if err != nil {
		return
	}

	r.msg = msg
}
