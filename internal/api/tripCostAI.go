package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func SendToAi(prompt string) (string, error) {

	apiKey := os.Getenv("GROQ_API_KEY")
	log.Println("KEY:", apiKey[:6], "...")
	if apiKey == "" {
		log.Println("API key is empty")
	}

	reqBody := ChatRequest{
		Model: "meta-llama/llama-4-scout-17b-16e-instruct",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var respJSON ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&respJSON); err != nil {
		return "", err
	}

	if len(respJSON.Choices) > 0 {
		return respJSON.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no respomse from AI")
}
