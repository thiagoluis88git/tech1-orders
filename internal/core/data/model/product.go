package model

import "gorm.io/gorm"

const (
	CategorySnack    = "Lanche"
	CategoryBeverage = "Bebida"
	CategoryDesert   = "Sobremesa"
	CategoryToppings = "Acompanhamento"
	CategoryCombo    = "Combo"
)

type Product struct {
	gorm.Model
	Name         string `gorm:"unique"`
	Description  string
	Category     string
	Price        float64
	ProductImage []ProductImage
	ComboProduct []ComboProduct
}

type ProductImage struct {
	gorm.Model
	ProductID uint
	ImageUrl  string
}

type ComboProduct struct {
	gorm.Model
	ProductID      uint
	ComboProductID uint
}
