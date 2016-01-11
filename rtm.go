package slackbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// start does a rtm.start, and returns a websocket URL and user ID. The
// websocket URL can be used to initiate an RTM session.
func startRTM(token string) (url, id string, err error) {
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
