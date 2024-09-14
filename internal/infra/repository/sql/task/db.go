package task

import (
	"fmt"

	"github.com/samber/lo"
	"gorm.io/gorm"

	"github.com/rubberduckkk/ducker/internal/domain/task"
	"github.com/rubberduckkk/ducker/internal/domain/task/valueobj"
)

type sqlRepository struct {
	db *gorm.DB
}

type Config struct {
	DB *gorm.DB
}

type Option func(*Config)

func WithDB(db *gorm.DB) Option {
	return func(c *Config) {
		c.DB = db
	}
}

func NewRepository(opts ...Option) task.Repository {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return &sqlRepository{
		db: cfg.DB,
	}
}

func (s *sqlRepository) Create(task *task.Task) error {
	model := FromEntity(task)
	return s.db.Model(model).Create(model).Error
}

func (s *sqlRepository) Update(task *task.Task) error {
	return s.db.Save(FromEntity(task)).Error
}

func (s *sqlRepository) Get(id string) (*task.Task, error) {
	var t Task
	if err := s.db.Model(t).First(&t, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return t.ToEntity(), nil
}

func (s *sqlRepository) Remove(id string) error {
	return s.db.Delete(&Task{}, "id = ?", id).Error
}

func (s *sqlRepository) List(customerID string, cursor int64, batchSize int, order valueobj.QueryOrder) (res []*task.Task, nextCursor int64, err error) {
	op := lo.Ternary(order == valueobj.OrderASC, ">=", "<=")
	query := fmt.Sprintf("customer_id = ? AND created_at %v ?", op)
	if err := s.db.Where(query, customerID, cursor).Limit(batchSize).Find(&res).Error; err != nil {
		return nil, 0, err
	}
	nextCursor = lo.Ternary(len(res) > 0, res[len(res)-1].CreatedAt.Unix(), cursor)
	return res, nextCursor, nil
}
