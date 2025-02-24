package models

type Holiday struct {
	Id      uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Date    string `gorm:"unique" json:"date"`
	Remarks string
}
