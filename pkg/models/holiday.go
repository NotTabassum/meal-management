package models

type Holiday struct {
	Id      uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Date    string `gorm:"primaryKey"`
	Remarks string
}
