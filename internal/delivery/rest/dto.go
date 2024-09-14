package rest

type Order string

const (
	OrderDesc Order = "desc"
	OrderAsc  Order = "asc"
)

var validOrders = map[Order]bool{
	OrderDesc: true,
	OrderAsc:  true,
}

func (o Order) IsValid() bool {
	return validOrders[o]
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}
