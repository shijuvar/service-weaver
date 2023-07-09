package paymentservice

import (
	"context"

	"github.com/ServiceWeaver/weaver"

	"github.com/shijuvar/service-weaver/orderapp/model"
)

type Service interface {
	MakePayment(ctx context.Context, orderPayment model.OrderPayment) error
}

type implementation struct {
	weaver.Implements[Service]
}

func (s *implementation) MakePayment(ctx context.Context, orderPayment model.OrderPayment) error {
	defer s.Logger().Info(
		"payment has been processed",
		"order ID", orderPayment.OrderID,
		"customer ID", orderPayment.CustomerID,
		"amount", orderPayment.Amount,
	)
	// ToDO: make the payment
	return nil
}
