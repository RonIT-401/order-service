package cgrpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	ccatalog "github.com/RonIT-401/order-service/internal/app/client/catalog"
	"github.com/RonIT-401/order-service/internal/app/entity"
	catalogv1 "github.com/RonIT-401/order-service/internal/pkg/grpc/gen/catalog/v1"
)

type client struct {
	raw catalogv1.CatalogServiceClient
}

var grpcCodeToAppError = map[codes.Code]error{
	codes.NotFound:        entity.ErrNotFound,
	codes.InvalidArgument: entity.ErrIncorrectParameters,
	codes.AlreadyExists:   entity.ErrAlreadyExists,
}

func normalizeError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return entity.ErrInternal
	}

	appErr, ok := grpcCodeToAppError[st.Code()]
	if ok {
		return appErr
	}

	return entity.ErrInternal
}

var _ ccatalog.Client = (*client)(nil)

func NewClient(address string) (ccatalog.Client, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	return &client{
		raw: catalogv1.NewCatalogServiceClient(conn),
	}, conn, nil
}

func (c *client) Ping(ctx context.Context) error {
	_, err := c.raw.GetProducts(ctx, &catalogv1.GetProductsRequest{})
	return err
}

func (c *client) GetProduct(ctx context.Context, req *catalogv1.GetProductRequest) (*catalogv1.GetProductResponse, error) {
	resp, err := c.raw.GetProduct(ctx, req)
	if err != nil {
		return nil, normalizeError(err)
	}

	return resp, nil
}

func (c *client) GetProducts(ctx context.Context, req *catalogv1.GetProductsRequest) (*catalogv1.GetProductsResponse, error) {
	resp, err := c.raw.GetProducts(ctx, req)
	if err != nil {
		return nil, normalizeError(err)
	}

	return resp, nil
}
