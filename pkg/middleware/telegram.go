package middleware

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"meal-management/pkg/config"
)

func SendTelegramMessage(message string) error {

	botToken := config.LocalConfig.TELEGRAM_BOT_TOKEN
	chatIDStr := config.LocalConfig.TELEGRAM_CHAT_ID

	if botToken == "" || chatIDStr == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN or TELEGRAM_CHAT_ID is missing in .env")
	}

	// Debug print
	fmt.Printf("Sending message to Telegram ChatID: %d\n", chatIDStr)
	fmt.Printf("Message: %s\n", message)

	// Create Telegram API URL
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	fmt.Println(url)
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"chat_id": chatIDStr,
			"text":    message,
		}).
		Post(url)

	if err != nil {
		return fmt.Errorf("failed to send POST request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("telegram API error: %s", resp.String())
	}

	fmt.Println("✅ Telegram message sent successfully!")
	return nil
}
