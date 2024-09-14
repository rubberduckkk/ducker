package customer

import (
	"github.com/rubberduckkk/ducker/internal/domain/customer"
	"github.com/rubberduckkk/ducker/internal/domain/customer/entity"
	"github.com/rubberduckkk/ducker/internal/domain/customer/valueobj"
)

type Customer struct {
	ID       string `json:"id" gorm:"primaryKey;not null"`
	Name     string `json:"name" gorm:"column:name;not null"`
	AreaCode string `json:"area_code" gorm:"column:area_code;not null"`
	PhoneNum string `json:"phone_num" gorm:"column:phone_num;not null"`
}

func (Customer) TableName() string {
	return "customers"
}

func FromEntity(customer *customer.Customer) *Customer {
	m := new(Customer)
	m.ID = customer.ID
	m.Name = customer.Name
	m.AreaCode = customer.AreaCode
	m.PhoneNum = customer.PhoneNum
	return m
}

func (m Customer) ToEntity() *customer.Customer {
	return &customer.Customer{
		Person: &entity.Person{
			ID:   m.ID,
			Name: m.Name,
		},
		ContactInfo: valueobj.ContactInfo{AreaCode: m.AreaCode, PhoneNum: m.PhoneNum},
	}
}
