package product

type IProductRepository interface {
	FindAll()([]*Product, error)
	FindById(id int64) (*Product, error)
	Create(product *Product) error
	Update(product *Product) error
	Delete(id int64) error	
}

type ICategoryRepository interface {
	FindAll()([]*Category, error)
	FindById(id int64)(*Category, error)
	Create(category *Category) error
	Update(category *Category) error
	Delete(id int64) error
}