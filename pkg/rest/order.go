package rest

type QueryOrder string

const (
	OrderDesc QueryOrder = "desc"
	OrderAsc  QueryOrder = "asc"
)

var validOrders = map[QueryOrder]bool{
	OrderDesc: true,
	OrderAsc:  true,
}

func (o QueryOrder) IsValid() bool {
	return validOrders[o]
}
