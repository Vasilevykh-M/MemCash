package database

import (
	"context"
	"gitlab.ozon.dev/go/classroom-9/students/homework-7/pkg/database"
	"time"
)

type Client struct {
	conn *database.Database
}

func NewClient(db *database.Database) *Client {
	return &Client{db}
}

func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	err := c.conn.Begin(ctx)
	if err != nil {
		return err
	}
	err = c.conn.Exec(ctx, "INSERT INTO table1 (a, b) VALUES (\"a\", \"b\")", value, key)
	if err != nil {
		err = c.conn.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}
	err = c.conn.Commit(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Get(ctx context.Context, key string) (any, error) {

	err := c.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	row := c.conn.QueryRow(ctx, "SELECT a FROM table1 WHERE b == a", key)

	err = c.conn.Commit(ctx)

	if err != nil {
		return nil, err
	}

	return row, nil
}
