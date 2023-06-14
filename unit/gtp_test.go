package unit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
)

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      string                   `json:"id"`
	Object  string                   `json:"object"`
	Created int                      `json:"created"`
	Model   string                   `json:"model"`
	Choices []map[string]interface{} `json:"choices"`
	Usage   map[string]interface{}   `json:"usage"`
}

type ChoiceItem struct {
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float32 `json:"temperature"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
	//ModelType        string  `json:"modeltype"`
}

const BASEURL = "https://api.openai.com/v1/"

func TestCompletions(t *testing.T) {
	msg := "你是chatgtp还是gtp-3"
	requestBody := ChatGPTRequestBody{
		Model:            "text-davinci-003",
		Prompt:           msg,
		MaxTokens:        2048,
		Temperature:      1,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		//ModelType:        "chat",
	}
	requestData, err := json.Marshal(requestBody)
	if err != nil {
		panic(err.Error())
	}
	log.Printf("request gtp json string : %v", string(requestData))

	req, err := http.NewRequest("POST", BASEURL+"completions", bytes.NewBuffer(requestData))
	if err != nil {
		panic(err.Error())
	}

	apiKey := "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"
	apiKey = "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"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	proxy := "http://127.0.0.1:7890/"
	proxyAddress, _ := url.Parse(proxy)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyAddress),
		},
	}
	response, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err.Error())
	}

	gptResponseBody := &ChatGPTResponseBody{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		panic(err.Error())
	}
	var reply string
	if len(gptResponseBody.Choices) > 0 {
		for _, v := range gptResponseBody.Choices {
			reply = v["text"].(string)
			break
		}
	}
	log.Printf("gpt response text: %s \n", reply)
	fmt.Println(reply)
}

// RequestData 结构体定义了向ChatGPT发送请求所需的参数
type RequestData struct {
	Model   string        `json:"model"` // 使用的模型
	Message []MessageData `json:"messages"`
}

type MessageData struct {
	Role    string `json:"role"` // 使用的模型
	Content string `json:"content"`
}

// ChatGPTResponseBody 请求体
type ChatGPTResponseBodyGtp struct {
	ID      string                   `json:"id"`
	Object  string                   `json:"object"`
	Created int                      `json:"created"`
	Model   string                   `json:"model"`
	Choices []map[string]interface{} `json:"choices"`
	Usage   map[string]interface{}   `json:"usage"`
}

// ResponseData 结构体定义了从ChatGPT返回的响应数据
type ResponseData struct {
	Text string `json:"text"` // 生成的文本
}

func TestGtp35(t *testing.T) {
	requestBody := RequestData{
		Model:   "gpt-3.5-turbo",
		Message: []MessageData{{Role: "user", Content: "你跟gtp-3的区别是什么"}},
	}
	requestData, err := json.Marshal(requestBody)
	if err != nil {
		panic(err.Error())
	}
	log.Printf("request gtp json string : %v", string(requestData))

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestData))
	if err != nil {
		panic(err.Error())
	}

	//apiKey := ""
	apiKey := ""

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	proxy := "http://127.0.0.1:7890/"
	proxyAddress, _ := url.Parse(proxy)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyAddress),
		},
	}
	response, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err.Error())
	}

	gptResponseBody := &ChatGPTResponseBody{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		panic(err.Error())
	}
	var reply string
	if len(gptResponseBody.Choices) > 0 {
		for _, v := range gptResponseBody.Choices {
			message, ok := v["message"].(map[string]interface{})
			if ok {
				reply = message["content"].(string)
				break
			}
		}
	}
	log.Printf("gpt response text: %s \n", reply)
	fmt.Println(reply)
}
