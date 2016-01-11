package slackbot

import (
	"log"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

var counter uint64 // atomic counter used to increment message IDs

type Message struct {
	Id      uint64 `json:"id"`
	Type    string `json:"type"`
	Subtype string `json:"subtype"`
	Channel string `json:"channel"`
	User    string `json:"user"`
	Text    string `json:"text"`
}

type Bot struct {
	id string
	ws *websocket.Conn
}

// New initializes a new bot session, returning a *Bot.
func New(token string) *Bot {
	url, id, err := startRTM(token)
	if err != nil {
		log.Fatal(err)
	}

	ws, err := websocket.Dial(url, "", "https://api.slack.com/")
	if err != nil {
		log.Fatal(err)
	}

	return &Bot{id, ws}
}

func (bot *Bot) Read() (msg Message, err error) {
	err = websocket.JSON.Receive(bot.ws, &msg)
	return msg, err
}

func (bot *Bot) Post(msg Message) error {
	msg.Id = atomic.AddUint64(&counter, 1)
	return websocket.JSON.Send(bot.ws, msg)
}
