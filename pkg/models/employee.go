package models

import "gorm.io/datatypes"

type Employee struct {
	EmployeeId     uint   `gorm:"primaryKey;autoIncrement" json:"employee_id"`
	Name           string `json:"name"`
	Email          string `gorm:"unique; not null"`
	PhoneNumber    string `gorm:"type:char(11)" validate:"len=11,numeric" json:"phone_number"`
	Password       string
	DeptID         int
	Remarks        string
	DefaultStatus  bool
	StatusUpdated  bool `json:"status_updated"`
	IsAdmin        bool
	Photo          string         `json:"photo"`
	PreferenceFood datatypes.JSON `gorm:"type:json"`
}
