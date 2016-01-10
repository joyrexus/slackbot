package slackbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

// NewSocket initiates a websocket-based RTM API session. It returns 
// the websocket and the ID of the (bot-)user whom the token belongs to.
func NewSocket(token string) (*websocket.Conn, string) {
	wsurl, id, err := start(token)
	if err != nil {
		log.Fatal(err)
	}

	ws, err := websocket.Dial(wsurl, "", "https://api.slack.com/")
	if err != nil {
		log.Fatal(err)
	}

	return ws, id
}

// start does a rtm.start, and returns a websocket URL and user ID. The
// websocket URL can be used to initiate an RTM session.
func start(token string) (url, id string, err error) {
	url = fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", token)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("API request failed with code %d", resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", "", err
	}

	var data struct {
		Ok    bool
		Error string
		Url   string
		Self  struct {
			Id string
		}
	}
	if err = json.Unmarshal(body, &data); err != nil {
		return "", "", err
	}
	if !data.Ok {
		err = fmt.Errorf("Slack error: %s", data.Error)
		return "", "", err
	}

	return data.Url, data.Self.Id, nil
}

type Message struct {
	Id      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func read(ws *websocket.Conn) (m Message, err error) {
	err = websocket.JSON.Receive(ws, &m)
	return m, err
}

var counter uint64

func post(ws *websocket.Conn, m Message) error {
	m.Id = atomic.AddUint64(&counter, 1)
	return websocket.JSON.Send(ws, m)
}
