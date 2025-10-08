package product

type Product struct {
	ID int64
	Name string
	Description string
	Price float64
	CategoryId []int64
}

type Category struct {
	ID int64
	Name string
}