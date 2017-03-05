package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getLanguage(text string) string {
	data := map[string]string{"textIn": text}
	dataBytes, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "https://jmlk74oovf.execute-api.eu-west-1.amazonaws.com/dev/language?wait=true", bytes.NewBuffer(dataBytes))
	req.Header.Set("x-api-key", "9CAfxmC4WB10tnS9RY9oG92Io0M4trVp7HpTUEjR")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	fmt.Println(text)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	jsonBody := map[string]interface{}{}
	json.Unmarshal(body, &jsonBody)
	return jsonBody["results"].(map[string]interface{})["language"].(string)
}

func getSentiment(text, language string) string {
	data := map[string]string{"textIn": text, "language": language}
	dataBytes, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "https://jmlk74oovf.execute-api.eu-west-1.amazonaws.com/dev/sentiment?wait=true", bytes.NewBuffer(dataBytes))
	req.Header.Set("x-api-key", "9CAfxmC4WB10tnS9RY9oG92Io0M4trVp7HpTUEjR")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	jsonBody := map[string]interface{}{}
	json.Unmarshal(body, &jsonBody)
	return jsonBody["results"].(map[string]interface{})["prediction"].(string)
}
