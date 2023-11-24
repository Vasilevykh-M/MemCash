package cache

import (
	"context"
	"gitlab.ozon.dev/go/classroom-9/students/homework-7/internal/datasource"
	"gitlab.ozon.dev/go/classroom-9/students/homework-7/pkg/cashe"
	"gitlab.ozon.dev/go/classroom-9/students/homework-7/pkg/database"
	"time"
)

type Client struct {
	conn    *cashe.Cache
	source  datasource.Datasource
	durCash time.Duration
}

func NewClient(durationCleanUp, durCash time.Duration, db datasource.Datasource) *Client {
	return &Client{cashe.New(durationCleanUp), db, durCash}
}

func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	err := c.source.Set(ctx, key, value, expiration)
	if err != nil {
		return err
	}
	c.conn.Set(ctx, key, value, expiration)
	return nil
}

func (c *Client) Get(ctx context.Context, key string) (any, error) {
	value, err := c.conn.Get(ctx, key)
	if err != nil {
		value, err := c.source.Get(ctx, key)
		if err != nil {
			return nil, err
		}

		c.conn.Set(ctx, key, value, c.durCash)

		return value.(database.Row).Values[0], nil
	}

	return value.(string), nil
}
