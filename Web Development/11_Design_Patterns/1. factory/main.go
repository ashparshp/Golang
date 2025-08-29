package main

import (
	"fmt"
	"myapp/products"
)

func main() {
	// Create a variable of the type products.Product.
	factory := products.Product{}

	// Call the New() method on the variable, to get back a Product
	// with sensible defaults set.
	product := factory.New()

	// Print out the results.
	fmt.Println("My product was created at", product.CreatedAt.UTC())

	// The above logic is functionally identical to doing it by hand, as per the
	// commented out code below:
	// product := products.Product {
	// 	ProductName: "widget",
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(),
	// }
	// fmt.Println("My product was created at", product.CreatedAt.UTC())
}
