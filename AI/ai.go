package AI

import (
	"DiscordBot/cmd"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Promt - функция для общения с ии. Мы передаем промт и получаем ответ. Включает параметы:
// user - тот кто пишет
// promt - сообщение для ии
// sysPromt - системное сообщение для бота
// api - апи ключ для бота
func Promt(user, promt, sysPromt, api string, ratelimiter *cmd.SimpleRateLimiter) (string, error) {
	_, ok := ratelimiter.CheckLimit()
	if !ok {
		return user + " не нужно на меня так наседать! Я не скорострел.", nil
	}
	ratelimiter.Unlock(user)

	UserPromt := "Тебе написал " + user + ": " + promt
	var response map[string]interface{}
	url := "https://api.intelligence.io.solutions/api/v1/chat/completions"

	// Создаем тело запроса (пример)
	payload := strings.NewReader(fmt.Sprintf(`{
		"model": "meta-llama/Llama-3.3-70B-Instruct",
		"messages": [
			{"role": "system", "content": "%s"},
			{"role": "user", "content": "%s"}
		],
		"temperature": 0.7,
		"max_tokens": 500
	}`, sysPromt, UserPromt))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "error", err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+api)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "error", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "error", err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "Не хочу тебе отвечать, динаху", err
	}
	content := response["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	return content, nil
}
