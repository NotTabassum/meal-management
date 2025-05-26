package models

import "gorm.io/datatypes"

type Employee struct {
	EmployeeId uint   `gorm:"primaryKey;autoIncrement" json:"employee_id"`
	Name       string `json:"name"`
	Email      string `gorm:"unique" json:"email,omitempty"`
	Password   string
	DeptID     int
	Remarks    string
	//DefaultStatus *bool
	DefaultStatusLunch  *bool
	DefaultStatusSnacks *bool
	StatusUpdated       bool `json:"status_updated"`
	IsAdmin             bool
	PhoneNumber         string         `gorm:"type:char(11);unique;not null" validate:"len=11,numeric" json:"phone_number"`
	Photo               string         `json:"photo"`
	PreferenceFood      datatypes.JSON `gorm:"type:json"`
	IsPermanent         *bool          `gorm:"default:true" json:"is_permanent"`
	IsActive            *bool          `gorm:"default:true" json:"is_active"`
	Roll                string         `gorm:"default:''" json:"roll"`
	Designation         string         `gorm:"default:''" json:"designation"`
}
