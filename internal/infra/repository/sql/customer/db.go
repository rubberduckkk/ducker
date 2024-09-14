package customer

import (
	"gorm.io/gorm"

	"github.com/rubberduckkk/ducker/internal/domain/customer"
)

type sqlRepository struct {
	db *gorm.DB
}

type Config struct {
	DB *gorm.DB
}

type Option func(*Config)

func WithDB(db *gorm.DB) Option {
	return func(cfg *Config) {
		cfg.DB = db
	}
}

func NewRepository(opts ...Option) customer.Repository {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return &sqlRepository{
		db: cfg.DB,
	}
}

func (s *sqlRepository) Create(customer *customer.Customer) error {
	model := FromEntity(customer)
	return s.db.Model(model).Create(model).Error
}

func (s *sqlRepository) Update(customer *customer.Customer) error {
	return s.db.Save(FromEntity(customer)).Error
}

func (s *sqlRepository) Get(id string) (*customer.Customer, error) {
	var c Customer
	if err := s.db.Model(&c).Where("id = ?", id).First(&c).Error; err != nil {
		return nil, err
	}
	return c.ToEntity(), nil
}

func (s *sqlRepository) Remove(id string) error {
	return s.db.Delete(&Customer{}, id).Error
}
