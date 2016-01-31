package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	slackURL = "https://slack.com/api/%s"
)

// Call calls a Slack API method, setting the token of bot in the method call
// parameters.
func (bot *Bot) Call(method string, data url.Values) (map[string]interface{}, error) {
	data.Set("token", bot.Token)
	return callAPI(method, data)
}

func callAPI(method string, data url.Values) (map[string]interface{}, error) {
	methodURL := fmt.Sprintf(slackURL, method)
	response, err := http.PostForm(methodURL, data)
	return httpToJSON(response, err)
}

func httpToJSON(response *http.Response, err error) (map[string]interface{}, error) {
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return unpackJSON(body)
}

func unpackJSON(body []byte) (map[string]interface{}, error) {
	var payload map[string]interface{}
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
