package llm

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type OpenAIClient struct {
	apiKey string
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{apiKey: apiKey}
}

func (c *OpenAIClient) Call(prompt string) (string, error) {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"model":      "gpt-3.5-turbo",
		"messages":   []map[string]string{{"role": "user", "content": prompt}},
		"max_tokens": 1000,
	})

	resp, err := http.Post(
		"https://api.openai.com/v1/chat/completions",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string), nil
}