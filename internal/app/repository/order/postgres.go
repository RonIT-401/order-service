package porder

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"

	"github.com/RonIT-401/order-service/internal/app/entity"
	"github.com/RonIT-401/order-service/internal/app/repository"
	rcpostgres "github.com/RonIT-401/order-service/internal/app/repository/conn/postgres"
)

type repoPg struct {
	conn *rcpostgres.Client
}

func NewRepo(client *rcpostgres.Client) repository.Order {
	return &repoPg{conn: client}
}

func (r *repoPg) InsideTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.conn.InsideTx(ctx, fn)
}

func (r *repoPg) Create(ctx context.Context, order entity.Order) error {
	return r.conn.GetDB(ctx).
		WithContext(ctx).
		Create(&order).
		Error
}

func (r *repoPg) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Order, error) {
	var order entity.Order

	err := r.conn.GetDB(ctx).
		WithContext(ctx).
		Preload("Items").
		Where("guid = ?", guid).
		First(&order).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Order{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.Order{}, err
	}

	return order, nil
}

func (r *repoPg) Update(ctx context.Context, order entity.Order) error {
	result := r.conn.GetDB(ctx).
		WithContext(ctx).
		Model(&entity.Order{}).
		Where("guid = ?", order.GUID).
		Updates(map[string]any{
			"status":     order.Status,
			"updated_at": order.UpdatedAt,
		})

	if result.RowsAffected == 0 {
		return entity.ErrNotFound
	}

	return result.Error
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) error {
	result := r.conn.GetDB(ctx).
		WithContext(ctx).
		Where("guid = ?", guid).
		Delete(&entity.Order{})

	if result.RowsAffected == 0 {
		return entity.ErrNotFound
	}

	return nil
}

func (r *repoPg) List(ctx context.Context, status *string, userGUID *uuid.UUID) ([]entity.Order, error) {
	var orders []entity.Order

	query := r.conn.GetDB(ctx).WithContext(ctx)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if userGUID != nil {
		query = query.Where("user_guid = ?", *userGUID)
	}

	err := query.Find(&orders).Error

	return orders, err
}
