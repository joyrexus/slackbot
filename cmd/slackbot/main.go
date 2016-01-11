package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joyrexus/slackbot"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: slackbot TOKEN\n")
		os.Exit(1)
	}

	// start a websocket-based RTM API session
	bot := slackbot.New(os.Args[1])
	fmt.Println("bot is ready, use ^C to quit")

	stopword := "EOF"
	record := false
	queue := make(map[string][]string)

	for { 
		msg, err := bot.Read() // read each incoming event
		if err != nil {
			log.Fatal(err)
		}

		switch {

		// filter anything that's not a basic message event
		case msg.Type != "message" || msg.Subtype != "":	
			continue 

		// look for `<<` prefix
		case strings.HasPrefix(msg.Text, "&lt;&lt;"):
			record = true
			stopword = strings.TrimPrefix(msg.Text, "&lt;&lt;")
			msg.Text = fmt.Sprintf("heredoc started; stop with %q", stopword)
			if err := bot.Post(msg); err != nil {
				fmt.Println(err)
			}

		// look for stopword
		case msg.Text == stopword:
			msg.Text = "heredoc stopped!"
			if err := bot.Post(msg); err != nil {
				fmt.Println(err)
			}

		// capture messages if recording	
		case record == true:
			item := fmt.Sprintf("%s: %s", msg.User, msg.Text)
			queue[msg.Channel] = append(queue[msg.Channel], item)
			fmt.Println(queue)
		}
	}
}
