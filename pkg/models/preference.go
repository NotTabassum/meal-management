package models

type Preference struct {
	FoodId int    `gorm:"primary_key;AUTO_INCREMENT" json:"food_Id"`
	Food   string `json:"food"`
}
