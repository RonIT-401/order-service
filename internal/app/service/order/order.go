package sorder

import (
	"context"
	"time"

	"github.com/gofrs/uuid"

	"github.com/RonIT-401/order-service/internal/app/entity"
	"github.com/RonIT-401/order-service/internal/app/repository"
	"github.com/RonIT-401/order-service/internal/app/service"
)

type srv struct {
	repoOrder repository.Order
}

func NewService(repoOrder repository.Order) service.Order {
	return &srv{repoOrder: repoOrder}
}

func (s *srv) Create(ctx context.Context, req entity.RequestOrderCreate) (entity.Order, error) {
	now := time.Now()
	var totalPrice int64

	items := make([]entity.OrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		itemGUID := uuid.Must(uuid.NewV4())

		totalPrice += int64(item.Quantity) * item.UnitPrice

		items = append(items, entity.OrderItem{
			GUID:        itemGUID,
			ProductGUID: item.ProductGUID,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}

	order := entity.Order{
		GUID:       uuid.Must(uuid.NewV4()),
		UserGUID:   req.UserGUID,
		TotalPrice: totalPrice,
		Currency:   req.Currency,
		Status:     "pending",
		CreatedAt:  now,
		UpdatedAt:  now,
		Items:      items,
	}

	if err := s.repoOrder.Create(ctx, order); err != nil {
		return entity.Order{}, err
	}

	return order, nil
}

func (s *srv) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Order, error) {
	return s.repoOrder.GetByGUID(ctx, guid)
}

func (s *srv) Update(ctx context.Context, guid uuid.UUID, req entity.RequestOrderUpdate) (entity.Order, error) {
	var order entity.Order

	err := s.repoOrder.InsideTx(ctx, func(ctx context.Context) error {
		var err error

		order, err = s.GetByGUID(ctx, guid)
		if err != nil {
			return err
		}

		order.Status = req.Status

		return s.repoOrder.Update(ctx, order)
	})
	if err != nil {
		return entity.Order{}, err
	}

	return order, nil
}

func (s *srv) Delete(ctx context.Context, guid uuid.UUID) error {
	err := s.repoOrder.InsideTx(ctx, func(ctx context.Context) error {
		_, err := s.GetByGUID(ctx, guid)
		if err != nil {
			return err
		}

		return s.repoOrder.Delete(ctx, guid)
	})

	return err
}

func (s *srv) List(ctx context.Context, req entity.RequestOrderList) ([]entity.Order, error) {
	return s.repoOrder.List(ctx, req.Status, req.UserGUID)
}
