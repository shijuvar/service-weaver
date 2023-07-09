package notificationservice

import (
	"context"

	"github.com/ServiceWeaver/weaver"

	"github.com/shijuvar/service-weaver/orderapp/model"
)

type Service interface {
	Send(ctx context.Context, notification model.Notification) error
}

type implementation struct {
	weaver.Implements[Service]
}

func (s *implementation) Send(ctx context.Context, notification model.Notification) error {
	defer s.Logger().Info(
		"notification has been sent",
		"order ID", notification.OrderID,
		"customer ID", notification.CustomerID,
		"event", notification.Event,
		"modes", notification.Modes,
	)
	// ToDO: send the notification
	return nil
}
