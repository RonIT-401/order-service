package rcpostgres

import (
	"context"

	"gorm.io/gorm"
)

type contextKeyTx struct{}

func getTxFromContext(ctx context.Context) *gorm.DB {
	value, ok := ctx.Value(contextKeyTx{}).(*gorm.DB)
	if !ok {
		return nil
	}

	return value
}

func ctxWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, contextKeyTx{}, tx)
}
