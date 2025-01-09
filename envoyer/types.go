package envoyer

type ResponseBody struct {
	Status     string `json:"status"`
	Data       string `json:"data"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type TemplateVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type IndividualMessage struct {
	Receiver          string             `json:"receiver"`
	TemplateVariables []TemplateVariable `json:"variables"`
}
