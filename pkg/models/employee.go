package models

type Employee struct {
	EmployeeId    uint `gorm:"primaryKey;autoIncrement"`
	Name          string
	Email         string `gorm:"unique; not null"`
	Password      string
	DeptID        int
	Remarks       string
	DefaultStatus bool
	IsAdmin       bool
}
