package products

import "time"

// Product is the type for all products.
type Product struct {
	ProductName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// New is the factory pattern for this module. It creates a pointer
// to the type Product, populates CreatedAt and UpdatedAt with
// sensible defaults, and returns it.
func (p *Product) New() *Product {
	product := Product{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &product
}
