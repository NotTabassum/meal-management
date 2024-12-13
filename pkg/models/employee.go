package models

type Employee struct {
	EmployeeId    uint `gorm:"primaryKey;autoIncrement"`
	Name          string
	Email         string
	Password      string
	DeptID        int
	Remarks       string
	DefaultStatus bool
}
