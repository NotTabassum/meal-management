package models

type Department struct {
	DeptID   uint `gorm:"primaryKey;autoIncrement"`
	DeptName string
	Weekend  string
}
