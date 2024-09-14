package entity

import (
	"github.com/google/uuid"
)

type Person struct {
	ID   string
	Name string
}

func NewPerson(name string) *Person {
	return &Person{
		ID:   uuid.NewString(),
		Name: name,
	}
}
