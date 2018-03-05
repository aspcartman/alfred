package torrent

import (
	"github.com/aspcartman/alfred/telegram"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"bytes"
	"github.com/aspcartman/exceptions"
	"fmt"
	"os"
	"gopkg.in/tucnak/telebot.v2"
	"net/http"
	"github.com/satori/go.uuid"
	"strings"
	"time"
	"github.com/aspcartman/alfred/env"
)

type Bot struct {
	addr   string
	client *torrent.Client
	files  map[string]*torrent.File
}

func NewBot(addr string) *Bot {
	client, err := torrent.NewClient(&torrent.Config{
		DataDir:    os.TempDir(),
		NoUpload:   true,
		ListenAddr: fmt.Sprintf(":%d", 12312),
	})
	if err != nil {
		e.Throw("failed creating torrent client", err)
	}

	b := &Bot{addr, client, map[string]*torrent.File{}}

	go func() { e.Must(http.ListenAndServe(addr, b), "failed listening and serving torrent bot") }()

	return b
}

func (b *Bot) Handle(s *telegram.Session, msg *telebot.Message) {
	s.Reply("Оп-па, торрент. Ща открою..")

	data := s.Download(msg.Document.FileID)
	info, err := metainfo.Load(bytes.NewReader(data))
	if err != nil {
		s.Reply("Что-то не открывается твой торрент =(")
		e.Throw("error opening torrent file", err)
	}

	t, err := b.client.AddTorrent(info)
	if err != nil {
		s.Reply("Эм, не смог его открыть =(")
		e.Throw("error adding torrent file", err)
	}

	files := t.Files()
	names := []string{}
	for _, f := range files {
		names = append(names, f.DisplayPath())
	}

	file := files[s.Ask("Че качаем?", names...)]
	res := s.Reply("Стартую закачку...")

	id := uuid.NewV4().String()
	b.files[id] = file
	url := "http://" + b.addr + "/" + "torrent" + "/" + id

	go b.renderProgressRoutine(s, file, res, url)
}

func (b *Bot) renderProgressRoutine(s *telegram.Session, f *torrent.File, msg *telegram.Response, url string) {
	defer e.Catch(func(e *e.Exception) {
		env.Log.Error("failure rendering torrent download progress")
		go b.renderProgressRoutine(s, f, msg, url)
	})

	for i := 0; ; i = (i + 1) % 4 {
		time.Sleep(5 * time.Second)

		name := f.DisplayPath()
		progress := float64(f.Torrent().BytesCompleted()) / float64(f.Torrent().Length()) * 100

		ind := []string{"|", "/", "-", "\\"}[i]

		state := fmt.Sprintf("Качаю %s\n%f%% %s\n%s\n", name, progress, ind, url)

		t := f.Torrent()
		start, stop := b.pieceRange(f)
		for c, i := 0, start; c < 10 && i <= stop; i++ {
			piece := t.Piece(i)

			switch  {
			case piece.Storage().Completion().Complete:
				continue
			case c==0:
				piece.SetPriority(torrent.PiecePriorityNow)
			case c==1:
				piece.SetPriority(torrent.PiecePriorityNext)
			case c==2:
				piece.SetPriority(torrent.PiecePriorityReadahead)
			default:
				piece.SetPriority(torrent.PiecePriorityHigh)
			}

			c++
		}

		for _, st := range f.State() {
			switch {
			case st.Complete:
				state += "|"
			case st.Partial:
				state += ":"
			default:
				state += "."
			}
		}

		msg.Edit(state)

		if progress > 99 {
			s.Reply("Done")
			return
		}
	}
}

func (b *Bot) pieceRange(f *torrent.File) (int, int) {
	off := f.Offset()
	ln := f.Length()
	tor := f.Torrent()
	pieceLen := tor.Info().PieceLength
	return int(off / pieceLen), int((off + ln) / pieceLen)
}

func (b *Bot) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	env.Log.WithField("url", r.RequestURI).Info("requesting torrent")

	if !strings.HasPrefix(r.RequestURI, "/torrent/") {
		rw.WriteHeader(404)
		env.Log.WithField("url", r.RequestURI).Error("torrent handler called, but url is wrong")
		return
	}

	id := r.RequestURI[len("/torrent/"):]
	file, ok := b.files[id]
	if !ok {
		rw.WriteHeader(404)
		env.Log.WithField("id", id).Error("torrent was not found")
		return
	}

	reader := file.NewReader()
	rw.Header().Set("Content-Disposition", "attachment; filename=\""+file.DisplayPath()+"\"")

	http.ServeContent(rw, r, file.DisplayPath(), time.Now(), reader)
}
