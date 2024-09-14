package customer

import (
	"github.com/rubberduckkk/ducker/internal/domain/customer/entity"
	"github.com/rubberduckkk/ducker/internal/domain/customer/valueobj"
)

type Customer struct {
	*entity.Person
	valueobj.ContactInfo
}

func NewCustomer(name string, contact valueobj.ContactInfo) *Customer {
	return &Customer{
		Person:      entity.NewPerson(name),
		ContactInfo: contact,
	}
}
