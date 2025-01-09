package envoyer

import "time"

const Email = "email"

type EmailReq struct {
	EventName    string             `json:"event_name"`
	DeliveryTime *time.Time         `json:"delivery_time"`
	Receivers    []string           `json:"receivers"`
	Variables    []TemplateVariable `json:"variables"`
	Sender       string             `json:"sender"`
	Cc           []string           `json:"cc"`
	Bcc          []string           `json:"bcc"`
	Language     string             `json:"language"`
}

type BulkEmailReq struct {
	EventName              string              `json:"event_name"`
	DeliveryTime           *time.Time          `json:"delivery_time"`
	ReceiversWithVariables []IndividualMessage `json:"receivers_with_variables"`
	Sender                 string              `json:"sender"`
	Language               string              `json:"language"`
}

type emailRequestBody struct {
	Envoyer
	EmailReq
}

type bulkEmailRequestBody struct {
	Envoyer
	BulkEmailReq
}

func (e Envoyer) SendEmail(body EmailReq) (ResponseBody, error) {
	req := &emailRequestBody{e, body}
	return e.sendNotification(req, Email)
}

func (e Envoyer) SendBulkEmail(body BulkEmailReq) (ResponseBody, error) {
	req := &bulkEmailRequestBody{e, body}
	return e.sendNotification(req, Email)
}
