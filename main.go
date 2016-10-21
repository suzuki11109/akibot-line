package main

import (
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/callback", callbackHandler(bot))

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

func callbackHandler(bot *linebot.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		events, err := bot.ParseRequest(r)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if err = handleTextMessage(bot, message, event.ReplyToken); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}
}

func handleTextMessage(bot *linebot.Client, message *linebot.TextMessage, replyToken string) error {
	var replyMessage string

	switch message.Text {
	case "สวัสดี":
		replyMessage = "ดีจ้า"
	}

	if replyMessage != "" {
		if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
			return err
		}
	}

	return nil
}
