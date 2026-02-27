package handler

import (
	"net/http"
	"github.com/SephirothGit/Backend-service/internal/service"
)

func NewOrderHandlerV2(svc service.OrderService) http.HandlerFunc {
	return NewOrderHandler(svc)
}