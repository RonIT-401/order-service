package eprocessor

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/RonIT-401/order-service/internal/app/entity"
	ehandler "github.com/RonIT-401/order-service/internal/app/handler/event"
	"github.com/RonIT-401/order-service/internal/app/processor"
	"github.com/RonIT-401/order-service/internal/pkg/broker"
)

type orderDeliveryCalculatedProc struct {
	h   ehandler.OrderDeliveryCalculated
	bus broker.Bus[entity.EventOrderDeliveryCalculated]
}

func NewOrderDeliveryCalculatedEventsCatcher(
	h ehandler.OrderDeliveryCalculated,
	bus broker.Bus[entity.EventOrderDeliveryCalculated],
) processor.Processor {
	return &orderDeliveryCalculatedProc{h, bus}
}

func (p *orderDeliveryCalculatedProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	if err := p.bus.Subscribe(ctx, wg, p.h.CallbackOrderDeliveryCalculated); err != nil {
		log.Fatal().Err(err).Msg("failed to subscribe to order delivery calculated topic")
	}

	log.Info().
		Str("topic_name", p.bus.QueueName()).
		Msg("subscribed to order delivery calculated topic")
}
