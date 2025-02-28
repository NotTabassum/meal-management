package models

type ExtraMeal struct {
	Date       string `json:"date"`
	LunchCount int    `json:"lunch_count"`
	SnackCount int    `json:"snack_count"`
}
