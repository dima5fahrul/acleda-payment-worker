package clients

import (
	"context"

	"gorm.io/gorm"
)

type YugabyteClient interface {
	Create(ctx context.Context, value interface{}) error
	First(ctx context.Context, dest interface{}, query interface{}, args ...interface{}) error
	Save(ctx context.Context, value interface{}) error
	GetDB() *gorm.DB
}

type yugabyteClient struct {
	db *gorm.DB
}

func NewYugabyteClient(db *gorm.DB) YugabyteClient {
	return &yugabyteClient{db: db}
}

func (c *yugabyteClient) Create(ctx context.Context, value interface{}) error {
	return c.db.WithContext(ctx).Create(value).Error
}

func (c *yugabyteClient) First(ctx context.Context, dest interface{}, query interface{}, args ...interface{}) error {
	return c.db.WithContext(ctx).Where(query, args...).First(dest).Error
}

func (c *yugabyteClient) Save(ctx context.Context, value interface{}) error {
	return c.db.WithContext(ctx).Save(value).Error
}

func (c *yugabyteClient) GetDB() *gorm.DB {
	return c.db
}
