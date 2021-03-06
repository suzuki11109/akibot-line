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
                    log.Println("userID:", event.Source.UserID)
					if err = handleTextMessage(bot, message, event.ReplyToken); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}
}

var messageEngine = map[string]string{
	"หวัดดี":   "ดีจ้า",
	"สวัสดี":   "ดีจ้า",
	"กินไร":    "ผัดเผ็ดปลาวาฬ",
	"ทราย":     "หนอน",
	"สัส":      "พูดเพราะๆสิ ปั๊ดโธ่!",
	"สัด":      "พูดเพราะๆสิ ปั๊ดโธ่!",
	"แจ่ม":     "แจ่ม แจแด็ม แจ่ม ว้าววว",
	"55555555": "ขำไร",
	"ถถถถถถถถ": "เปลี่ยนภาษาก่อนนะ",
	"เก่ง":     "สวดยอดไปเลยย",
	"ศพ":       "อย่าว่าเพื่อน",
	"สวย":      "สวยจริงหรอ?",
	"เที่ยว":   "ไปทะเลกันเถอะ",
	"สุรินทร์": "สุรินทร์ เป็นถิ่นมีหอยย~",
	"สาด":      "แสดดดดดด",
	"โอเค":     "โอเคจ้ะ",
	"ควย":      "ควยพ่อง",
	"ต๊ะ":      "โซโล่!!!",
	"หนอน":     "ดึ๊บๆ",
	"ขอบคุณ":   "แต๊งกิ้วนะ",
	"ขอบใจ":    "แต๊งกิ้วนะ",
	"ขอโทษ":    "ขยมโซ้มโต๊ก",
	"ขอโทด":    "ขยมโซ้มโต๊ก",
	"พ่อง":     "แม่ง",
	"เหี้ย":    "เอะอะก็เหี้ย อะไรก็เหี้ย เหี้ยผิดอะไร",
	"ตลก":      "ขำดิ",
	"อาย":      "ไอเหม่ง",
	"กินเหล้า": "โตแล้ว เลิกกินเหล้าได้แล้ว",
	"จวย":      "หัวเค",
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
