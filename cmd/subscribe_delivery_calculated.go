package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/RonIT-401/order-service/internal/app/builder"
)

func SubscribeDeliveryCalculated() *cli.Command {
	return &cli.Command{
		Name:            "subscribe-delivery-calculated",
		Aliases:         []string{"consume-delivery-calculated"},
		Usage:           "Запускает consumer топика order.delivery.calculated",
		Action:          cmdSubscribeDeliveryCalculated,
		HideHelpCommand: true,
	}
}

func cmdSubscribeDeliveryCalculated(cCtx *cli.Context) error {
	app := builder.NewBuilder(cCtx)
	app.BuildConfig()
	app.BuildMonitorOpenTelemetry()
	app.BuildBrokerKafka()
	app.BuildRepoConnPostgres()
	app.BuildRepoOrder()
	app.BuildServiceOrderDelivery()
	app.BuildBusOrderDeliveryCalculated()
	app.BuildProcHttp()

	app.Run()
	return nil
}
