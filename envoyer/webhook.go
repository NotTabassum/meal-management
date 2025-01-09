package envoyer

import "time"

const Webhook = "webhook"

type WebhookReq struct {
	DeliveryTime *time.Time  `json:"delivery_time"`
	Data         interface{} `json:"data"`
}

type webhookRequestBody struct {
	Envoyer
	WebhookReq
}

func (e Envoyer) SendWithWebhook(body WebhookReq) (ResponseBody, error) {
	req := &webhookRequestBody{e, body}
	return e.sendNotification(req, Webhook)
}
