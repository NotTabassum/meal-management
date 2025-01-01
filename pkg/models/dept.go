package models

import "gorm.io/datatypes"

type Department struct {
	DeptID   int            `gorm:"primaryKey" json:"dept_id"`
	DeptName string         `gorm:"not null" json:"dept_name"`
	Weekend  datatypes.JSON `gorm:"type:json" json:"weekend"`
}
