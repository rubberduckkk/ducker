package customer

type Repository interface {
	Create(customer *Customer) error
	Update(customer *Customer) error
	Get(id string) (*Customer, error)
	Remove(id string) error
}
