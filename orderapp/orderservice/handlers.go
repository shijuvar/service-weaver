package orderservice

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/shijuvar/service-weaver/orderapp/model"
)

var ctx = context.Background()

func (s *Server) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Invalid Order Data", 500)
		return
	}
	id, _ := uuid.NewUUID()
	aggregateID := id.String()
	order.ID = aggregateID
	order.Status = "Placed"
	order.CreatedOn = time.Now()
	order.Amount = order.GetAmount()
	s.Logger().Debug("items:", len(order.OrderItems))
	// persistence using orderRepository component
	//if err := s.orderRepository.Get().CreateOrder(ctx, order); err != nil {
	//	s.Logger().Error(
	//		"order failed",
	//		"error:", err,
	//	)
	//	http.Error(w, "order failed", http.StatusInternalServerError)
	//	return
	//}
	orderPayment := model.OrderPayment{
		OrderID:    order.ID,
		CustomerID: order.CustomerID,
		Amount:     order.Amount,
	}
	// make payment using paymentservice component
	if err := s.paymentService.Get().MakePayment(ctx, orderPayment); err != nil {
		s.Logger().Error(
			"payment failed",
			"error:", err,
		)
		http.Error(w, "payment failed", http.StatusInternalServerError)
		return
	}
	notification := model.Notification{
		OrderID:    order.ID,
		CustomerID: order.CustomerID,
		Event:      "order.placed",
		Modes:      []string{model.Email, model.SMS},
	}
	// send notification using notificationService component
	if err := s.notificationService.Get().Send(ctx, notification); err != nil {
		s.Logger().Error(
			"notification failed",
			"error:", err,
		)
	}
	s.Logger().Info(
		"order has been placed",
		"order id", order.ID,
		"customer id", order.CustomerID,
	)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderid")
	if orderID == "" {
		http.Error(w, "order ID is required", http.StatusBadRequest)
		return
	}
	// ToDO: Get Order by ID
	s.Logger().Debug(
		"GetOrderByID",
		"order ID", orderID,
	)
	w.WriteHeader(http.StatusOK)
}
