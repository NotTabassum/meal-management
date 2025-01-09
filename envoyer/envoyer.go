package envoyer

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Envoyer struct {
	url       string
	AppKey    string `json:"app_key"`
	ClientKey string `json:"client_key"`
}

func New(url string, appKey string, clientKey string) *Envoyer {
	return &Envoyer{
		url:       url,
		AppKey:    appKey,
		ClientKey: clientKey,
	}
}

func (e Envoyer) sendNotification(body interface{}, notificationType string) (ResponseBody, error) {
	requestBodyBytes, err := json.Marshal(body)
	if err != nil {
		return ResponseBody{}, err
	}
	return e.trigger(requestBodyBytes, notificationType)
}

func (e Envoyer) trigger(requestBodyBytes []byte, notificationType string) (ResponseBody, error) {
	req, err := http.NewRequest("POST", e.url+"/api/v2/publish/"+notificationType, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return ResponseBody{Message: "New request creation failed"}, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ResponseBody{Message: "Http request sent failed"}, err
	}
	defer resp.Body.Close()

	var responseBody ResponseBody
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return ResponseBody{StatusCode: resp.StatusCode, Message: "Response decode failed"}, err
	}

	responseBody.StatusCode = resp.StatusCode
	if responseBody.StatusCode != http.StatusOK {
		return responseBody, errors.New("failed to trigger the event")
	}
	responseBody.Message = "Message sent successfully"

	return responseBody, nil
}
