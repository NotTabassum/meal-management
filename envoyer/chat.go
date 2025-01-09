package envoyer

import "time"

const Chat = "Chat"

type ChatReq struct {
	EventName    string             `json:"event_name"`
	DeliveryTime *time.Time         `json:"delivery_time"`
	ProviderType string             `json:"provider_type"`
	Receivers    []string           `json:"receivers"`
	Variables    []TemplateVariable `json:"variables"`
	Language     string             `json:"language"`
}

type BulkChatReq struct {
	EventName              string              `json:"event_name"`
	DeliveryTime           *time.Time          `json:"delivery_time"`
	ReceiversWithVariables []IndividualMessage `json:"receivers_with_variables"`
	Language               string              `json:"language"`
}

type chatRequestBody struct {
	Envoyer
	ChatReq
}

type bulkChatRequestBody struct {
	Envoyer
	BulkChatReq
}

func (e Envoyer) SendChatMessage(body ChatReq) (ResponseBody, error) {
	req := &chatRequestBody{e, body}
	return e.sendNotification(req, Chat)
}

func (e Envoyer) SendBulkChatMessage(body BulkChatReq) (ResponseBody, error) {
	req := &bulkChatRequestBody{e, body}
	return e.sendNotification(req, Chat)
}
