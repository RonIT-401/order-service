package rcpostgres

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/RonIT-401/order-service/internal/app/config/section"
)

type Client struct {
	db  *gorm.DB
	cfg section.RepositoryPostgres
}

func (c *Client) DB() *gorm.DB {
	return c.db
}

func NewClient(ctx context.Context, cfg section.RepositoryPostgres) (*Client, error) {
	var u url.URL
	u.Scheme = "postgres"
	u.Host = cfg.Address
	u.User = url.UserPassword(cfg.Username, cfg.Password)
	u.Path = cfg.Name

	args := make(url.Values)
	args.Set("sslmode", "disable")
	u.RawQuery = args.Encode()

	dsn := u.String()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(10)

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return &Client{db: db, cfg: cfg}, nil
}

func (c *Client) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func (c *Client) GetDB(ctx context.Context) *gorm.DB {
	tx := getTxFromContext(ctx)
	if tx != nil {
		return tx
	}

	return c.db
}

func (c *Client) InsideTx(ctx context.Context, fn func(ctx context.Context) error) error {
	currentTx := getTxFromContext(ctx)

	if currentTx != nil {
		return fn(ctx)
	}

	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctxWithTx(ctx, tx))
	})
}
