package models

type Holiday struct {
	Date    string `gorm:"primaryKey"`
	Remarks string
}
