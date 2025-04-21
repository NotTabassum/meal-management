package middleware

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func SendTelegramMessage(message string) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")

	fmt.Println(botToken, chatIDStr)

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Fatal("Invalid TELEGRAM_CHAT_ID")
	}

	client := resty.New()
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	fmt.Println(message)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"chat_id": chatID,
			"text":    message,
		}).
		Post(url)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("Telegram API error: %s", resp.String())
	}

	return nil
}
