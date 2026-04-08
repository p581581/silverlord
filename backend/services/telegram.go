package services

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// SendTelegram 以 HTML parse mode 發送訊息至單一 Chat ID
func SendTelegram(botToken string, chatId int64, message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	params := url.Values{}
	params.Set("chat_id", fmt.Sprintf("%d", chatId))
	params.Set("text", message)
	params.Set("parse_mode", "HTML")

	resp, err := http.Post(apiURL, "application/x-www-form-urlencoded",
		strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status %d", resp.StatusCode)
	}
	return nil
}

// SendTelegramMulti 發送至逗號分隔的多組 Chat ID
func SendTelegramMulti(botToken, chatIds, message string) {
	for _, raw := range strings.Split(chatIds, ",") {
		id, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
		if err != nil || id == 0 {
			continue
		}
		SendTelegram(botToken, id, message)
	}
}
