package eorder

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/RonIT-401/order-service/internal/app/entity"
	ehandler "github.com/RonIT-401/order-service/internal/app/handler/event"
	"github.com/RonIT-401/order-service/internal/app/service"
	"github.com/RonIT-401/order-service/internal/pkg/broker"
)

type handler struct {
	order service.OrderDelivery
}

func NewHandler(
	order service.OrderDelivery,
) ehandler.OrderDeliveryCalculated {
	return &handler{order: order}
}

func (h *handler) CallbackOrderDeliveryCalculated(
	ctx context.Context, ev *entity.EventOrderDeliveryCalculated, _ []broker.Header,
) error {
	log.Info().
		Ctx(ctx).
		EmbedObject(ev).
		Msg("Получено событие order.delivery.calculated")

	if ev.OrderGUID == "" {
		return broker.NotCriticalError(fmt.Errorf("order.delivery.calculated: empty order_guid"))
	}

	applied, err := h.order.UpdateDelivery(ctx, *ev)
	if err != nil {
		return fmt.Errorf("order.delivery.calculated: %w", err)
	}

	if !applied {
		log.Warn().
			Ctx(ctx).
			Str("order_guid", ev.OrderGUID).
			Msg("Заказ не обновлён: статус отличается от pending или заказ не найден")

		return nil
	}

	log.Info().
		Ctx(ctx).
		Str("order_guid", ev.OrderGUID).
		Msg("Заказ обновлен")

	return nil
}
