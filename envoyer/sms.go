package envoyer

import "time"

const Sms = "sms"

type SmsReq struct {
	EventName    string             `json:"event_name"`
	DeliveryTime *time.Time         `json:"delivery_time"`
	Receivers    []string           `json:"receivers"`
	Variables    []TemplateVariable `json:"variables"`
	Language     string             `json:"language"`
}

type BulkSmsReq struct {
	EventName              string              `json:"event_name"`
	DeliveryTime           *time.Time          `json:"delivery_time"`
	ReceiversWithVariables []IndividualMessage `json:"receivers_with_variables"`
	Language               string              `json:"language"`
}

type smsRequestBody struct {
	Envoyer
	SmsReq
}

type bulkSmsRequestBody struct {
	Envoyer
	BulkSmsReq
}

func (e Envoyer) SendSms(body SmsReq) (ResponseBody, error) {
	req := &smsRequestBody{e, body}
	return e.sendNotification(req, Sms)
}

func (e Envoyer) SendBulkSms(body BulkSmsReq) (ResponseBody, error) {
	req := &bulkSmsRequestBody{e, body}
	return e.sendNotification(req, Sms)
}
