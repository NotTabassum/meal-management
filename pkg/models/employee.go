package models

type Employee struct {
	EmployeeId    uint   `gorm:"primaryKey;autoIncrement" json:"employee_id"`
	Name          string `json:"name"`
	Email         string `gorm:"unique; not null"`
	PhoneNumber   string `gorm:"type:char(11);unique;not null" validate:"len=11,numeric" json:"phone_number"`
	Password      string
	DeptID        int
	Remarks       string
	DefaultStatus bool
	IsAdmin       bool
	Photo         string `json:"photo"`
}
