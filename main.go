package main

import (
	"log"
	"net/http"
	"os"
	"strings"

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

var messageEngine = map[string]string{
	"กิ":     "อะกิไม่อยู่ คุยกับป๋มได้",
	"สวัสดี": "ดีจ้า",
	"กินไร":  "ไข่เจียวหมูสับ",
	"ทราย":   "หนอน",
	"สัด":    "พูดเพราะๆสิ ปั๊ดโธ่!",
	"แจ่ม":   "แจ่ม แจแด็ม แจ่ม ว้าววว",
	"55555":  "ขำไร",
}

func handleTextMessage(bot *linebot.Client, message *linebot.TextMessage, replyToken string) error {
	var replyMessage string

	replyMessage = reply(message.Text)

	if replyMessage != "" {
		if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
			return err
		}
	}

	return nil
}

func reply(message string) string {
	for k, v := range messageEngine {
		if strings.Contains(message, k) {
			return v
		}
	}

	return ""
}
