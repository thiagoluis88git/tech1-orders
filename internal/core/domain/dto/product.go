package dto

type ProductForm struct {
	Id               uint          `json:"id"`
	Name             string        `json:"name" validate:"required"`
	Description      string        `json:"description" validate:"required"`
	Category         string        `json:"category" validate:"required"`
	Price            float64       `json:"price" validate:"required"`
	Images           []ProducImage `json:"images" validate:"required"`
	ComboProductsIds *[]uint       `json:"comboProductsIds"`
}

type ProductResponse struct {
	Id            uint               `json:"id"`
	Name          string             `json:"name" validate:"required"`
	Description   string             `json:"description" validate:"required"`
	Category      string             `json:"category" validate:"required"`
	Price         float64            `json:"price" validate:"required"`
	Images        []ProducImage      `json:"images" validate:"required"`
	ComboProducts *[]ProductResponse `json:"comboProducts"`
}

type ProductCreationResponse struct {
	Id uint `json:"id"`
}

type ProducImage struct {
	ImageUrl string `json:"imageUrl" validate:"required"`
}

type ComboForm struct {
	Id          uint    `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	Products    []uint  `json:"products" validate:"required"`
}

type Combo struct {
	Id          uint          `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Price       float64       `json:"price"`
	Products    []ProductForm `json:"products"`
}
