package horder

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"

	"github.com/RonIT-401/order-service/internal/app/entity"
	rhandler "github.com/RonIT-401/order-service/internal/app/handler/http"
	"github.com/RonIT-401/order-service/internal/app/service"
	"github.com/RonIT-401/order-service/internal/pkg/http/httph"
)

type handler struct {
	srv service.Order
}

func NewHandler(srv service.Order) rhandler.Order {
	return &handler{srv: srv}
}

func mapOrderItems(items []entity.OrderItem) []entity.ResponseOrderItem {
	resp := make([]entity.ResponseOrderItem, 0, len(items))

	for _, item := range items {
		resp = append(resp, entity.ResponseOrderItem{
			GUID:        item.GUID,
			ProductGUID: item.ProductGUID,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		})
	}

	return resp
}

func (h *handler) Create(c *gin.Context) {
	var req entity.RequestOrderCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		httph.HandleError(c.Writer, c.Request, entity.ErrIncorrectParameters)
		return
	}

	order, err := h.srv.Create(c.Request.Context(), req)
	if err != nil {
		httph.HandleError(c.Writer, c.Request, err)
		return
	}

	resp := entity.ResponseOrderCreate{
		GUID:       order.GUID,
		UserGUID:   order.UserGUID,
		TotalPrice: order.TotalPrice,
		Currency:   order.Currency,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt,
		Items:      mapOrderItems(order.Items),
	}

	httph.SendJSON(c.Writer, http.StatusCreated, resp)
}

func (h *handler) GetByGUID(c *gin.Context) {
	guid, err := uuid.FromString(c.Param("guid"))
	if err != nil {
		httph.HandleError(c.Writer, c.Request, entity.ErrIncorrectParameters)
		return
	}

	order, err := h.srv.GetByGUID(c.Request.Context(), guid)
	if err != nil {
		httph.HandleError(c.Writer, c.Request, err)
		return
	}

	resp := entity.ResponseOrderGet{
		ID:         order.ID,
		GUID:       order.GUID,
		UserGUID:   order.UserGUID,
		TotalPrice: order.TotalPrice,
		Currency:   order.Currency,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		Items:      mapOrderItems(order.Items),
	}

	httph.SendJSON(c.Writer, http.StatusOK, resp)
}

func (h *handler) Update(c *gin.Context) {
	guid, err := uuid.FromString(c.Param("guid"))
	if err != nil {
		httph.HandleError(c.Writer, c.Request, entity.ErrIncorrectParameters)
		return
	}

	var req entity.RequestOrderUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		httph.HandleError(c.Writer, c.Request, entity.ErrIncorrectParameters)
		return
	}

	order, err := h.srv.Update(c.Request.Context(), guid, req)
	if err != nil {
		httph.HandleError(c.Writer, c.Request, err)
		return
	}

	resp := entity.ResponseOrderUpdate{
		ID:         order.ID,
		GUID:       order.GUID,
		UserGUID:   order.UserGUID,
		TotalPrice: order.TotalPrice,
		Currency:   order.Currency,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		Items:      mapOrderItems(order.Items),
	}

	httph.SendJSON(c.Writer, http.StatusOK, resp)
}

func (h *handler) Delete(c *gin.Context) {
	guid, err := uuid.FromString(c.Param("guid"))
	if err != nil {
		httph.HandleError(c.Writer, c.Request, entity.ErrIncorrectParameters)
		return
	}

	err = h.srv.Delete(c.Request.Context(), guid)
	if err != nil {
		httph.HandleError(c.Writer, c.Request, err)
		return
	}

	httph.SendEmpty(c.Writer, http.StatusNoContent)
}

func (h *handler) List(c *gin.Context) {
	var req entity.RequestOrderList

	if c.Request.ContentLength != 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			httph.HandleError(c.Writer, c.Request, entity.ErrIncorrectParameters)
			return
		}
	}

	orders, err := h.srv.List(c.Request.Context(), req)
	if err != nil {
		httph.HandleError(c.Writer, c.Request, err)
		return
	}

	resp := entity.ResponseOrderList{
		Data: make([]entity.ResponseOrderListItem, 0, len(orders)),
	}

	for _, order := range orders {
		resp.Data = append(resp.Data, entity.ResponseOrderListItem{
			GUID:       order.GUID,
			UserGUID:   order.UserGUID,
			TotalPrice: order.TotalPrice,
			Currency:   order.Currency,
			Status:     order.Status,
			CreatedAt:  order.CreatedAt,
			UpdatedAt:  order.UpdatedAt,
		})
	}

	httph.SendJSON(c.Writer, http.StatusOK, resp)
}
