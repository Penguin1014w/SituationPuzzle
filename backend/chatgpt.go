package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// DeepSeek API related constants
const (
	deepseekAPIURL = "https://api.deepseek.com/v1/chat/completions"
)

// DeepSeek API request struct
type DeepSeekRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// message struct
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DeepSeek API response struct
type DeepSeekResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// use AI to check the user's answer
func checkAnswerWithChatGPT(riddle Riddle, userAnswer string, lang string) (string, error) {
	// first do a simple string matching
	correctAnswer := riddle.AnswerCH
	content := riddle.ContentCH
	if lang == "EN" {
		correctAnswer = riddle.AnswerEN
		content = riddle.ContentEN
	}
	if strings.TrimSpace(userAnswer) == strings.TrimSpace(correctAnswer) {
		return "correct", nil
	}

	// build the AI prompt
	var prompt string
	if lang == "EN" {
		prompt = fmt.Sprintf(`This is a lateral thinking puzzle.
Riddle: %s
Correct Answer: %s
User's Answer: %s

Please carefully analyze the user's answer and return one of the following three states:
1. If the user's answer is logically consistent with the correct answer, return "correct"
2. If the user's answer is a leading question and the answer to that question is related to the correct answer, return "yes" or "no"
3. If the user's answer is completely irrelevant to the puzzle, return "irrelevant"

Criteria:
- Leading question: User is trying to get information through questions, such as "Is he...?", "Is this...?", etc.
- Related answer: The answer to the question is related to the key information of the puzzle
- Irrelevant answer: The question is unrelated to the scenario described in the riddle
- Correct answer: Must contain all key information points from the correct answer

Please return exactly in the following format:
- For leading questions that are related, return "yes" or "no"
- For correct answers, return "correct"
- For irrelevant questions, return "irrelevant"
Do not include any other text`, content, correctAnswer, userAnswer)
	} else {
		prompt = fmt.Sprintf(`这是一个海龟汤谜题。
谜面：%s
正确答案：%s
用户回答：%s

请仔细分析用户的回答，并返回以下三种状态之一：
1. 如果用户的回答与正确答案在核心逻辑上完全一致，返回 "correct"
2. 如果用户的回答是一个引导性问题，且问题的答案与正确答案相关，返回 "yes" 或 "no"
3. 如果用户的回答与本题完全无关，返回 "irrelevant"

判断标准：
- 引导性问题：用户试图通过提问来获取信息，如"他是...吗？"、"这是...吗？"等
- 相关回答：问题的答案与谜题的关键信息相关
- 无关回答：问题与谜面描述的场景无关，没有针对关键信息
- 正确回答：必须包含正确答案中的所有关键信息点

请严格按照以下格式返回：
- 如果是引导性问题且相关，返回 "yes" 或 "no"
- 如果是正确答案，返回 "correct"
- 如果是无关问题，返回 "irrelevant"
不要包含其他文字`, content, correctAnswer, userAnswer)
	}

	// build the API request
	requestBody := DeepSeekRequest{
		Model: "deepseek-chat",
		Messages: []Message{
			{
				Role: "system",
				Content: func() string {
					if lang == "EN" {
						return "You are a lateral thinking puzzle assistant. Please analyze if the user's answer is a leading question or a direct answer, and whether it's relevant to the puzzle. Your response must only contain 'yes', 'no', 'correct', or 'irrelevant', with no other text."
					}
					return "你是一个海龟汤谜题判断助手。请分析用户的回答是引导性问题还是直接答案，以及是否与题目相关。你的回答必须只包含 'yes'、'no'、'correct' 或 'irrelevant'，不要包含其他文字。"
				}(),
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// send the request to the DeepSeek API
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", deepseekAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("DEEPSEEK_API_KEY is not set in environment")
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// parse the API response
	var deepseekResp DeepSeekResponse
	if err := json.NewDecoder(resp.Body).Decode(&deepseekResp); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	if len(deepseekResp.Choices) == 0 {
		return "", fmt.Errorf("no response from DeepSeek")
	}

	// handle the AI's result
	result := deepseekResp.Choices[0].Message.Content
	fmt.Printf("Debug - Riddle: %s\n", content)
	fmt.Printf("Debug - Correct Answer: %s\n", correctAnswer)
	fmt.Printf("Debug - User Answer: %s\n", userAnswer)
	fmt.Printf("Debug - AI Response: %s\n", result)

	cleanedResult := strings.TrimSpace(strings.ToLower(result))
	if cleanedResult == "correct" || cleanedResult == "yes" || cleanedResult == "no" || cleanedResult == "irrelevant" {
		return cleanedResult, nil
	}

	// if cannot determine, default return irrelevant
	return "irrelevant", nil
}
