package product

import "ool/api/schema"

// Repository is a simple map that holds the mock product data, simulating a database.
// In a real application, this would be a structure containing database connection logic
// (e.g., a DynamoDB or Aurora client) and methods to interact with it.
var Repository = map[string]schema.Product{
	"10": {ID: "10", Name: "Chicken Waffle", Price: 12.99, Category: "Waffle"},
	"11": {ID: "11", Name: "Spicy Tofu Wrap", Price: 9.50, Category: "Wrap"},
	"12": {ID: "12", Name: "Classic Cheeseburger", Price: 11.00, Category: "Burger"},
	"13": {ID: "13", Name: "Sweet Potato Fries", Price: 4.00, Category: "Side"},
}
