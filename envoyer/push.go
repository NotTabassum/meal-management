package envoyer

import "time"

const Push = "push"

type PushReq struct {
	EventName    string             `json:"event_name"`
	DeliveryTime *time.Time         `json:"delivery_time"`
	Receivers    []string           `json:"receivers"`
	Variables    []TemplateVariable `json:"variables"`
	ImageUrl     string             `json:"image_url"`
	Data         map[string]string  `json:"data"`
	Topic        string             `json:"topic"`
	Condition    string             `json:"condition"`
	Language     string             `json:"language"`
}

type BulkPushReq struct {
	EventName              string              `json:"event_name"`
	DeliveryTime           *time.Time          `json:"delivery_time"`
	ReceiversWithVariables []IndividualMessage `json:"receivers_with_variables"`
	ImageUrl               string              `json:"image_url"`
	Data                   map[string]string   `json:"data"`
	Language               string              `json:"language"`
}

type pushRequestBody struct {
	Envoyer
	PushReq
}

type bulkPushRequestBody struct {
	Envoyer
	BulkPushReq
}

func (e Envoyer) SendPushNotification(body PushReq) (ResponseBody, error) {
	req := &pushRequestBody{e, body}
	return e.sendNotification(req, Push)
}

func (e Envoyer) SendBulkPushNotification(body BulkPushReq) (ResponseBody, error) {
	req := &bulkPushRequestBody{e, body}
	return e.sendNotification(req, Push)
}
