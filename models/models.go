package models

import (
	"gorm.io/gorm"
)

type Make struct {
	gorm.Model
	Name            string `json:"name" gorm:"text"`
	Foundation_Year int    `json:"foundation_year" gorm:"integer"`
	Cars            []Car  `json:"cars" gorm:"foreignkey:MakeID"`
}

type Car struct {
	gorm.Model
	Name   string `json:"name" gorm:"text"`
	MakeID uint   `json:"make_id" gorm:"integer"`
	Year   int    `json:"year" gorm:"integer"`
	Price  int    `json:"price" gorm:"integer"`
}

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"text"`
	Email    string `json:"email" gorm:"text"`
	Password string `json:"password" gorm:"text"`
}
