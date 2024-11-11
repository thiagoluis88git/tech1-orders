package database

import (
	"github.com/thiagoluis88git/tech1-orders/internal/core/data/model"

	"gorm.io/gorm"
)

type Database struct {
	Connection *gorm.DB
}

func ConfigDatabase(dialector gorm.Dialector) (*Database, error) {
	db, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		return &Database{}, err
	}

	db.AutoMigrate(
		&model.Order{},
		&model.OrderProduct{},
		&model.Product{},
		&model.ProductImage{},
		&model.ComboProduct{},
		&model.OrderTicketNumber{},
	)

	return &Database{
		Connection: db,
	}, nil
}
