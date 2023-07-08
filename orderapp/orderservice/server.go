package orderservice

import (
	"context"
	"net/http"
	notificationservice "service-weaver/orderapp/notificationservice"

	"github.com/ServiceWeaver/weaver"
	chi "github.com/go-chi/chi/v5"

	"service-weaver/orderapp/paymentservice"
)

type Server struct {
	weaver.Implements[weaver.Main]

	handler http.Handler

	paymentService      weaver.Ref[paymentservice.Service]
	notificationService weaver.Ref[notificationservice.Service]

	orderapi weaver.Listener `weaver:"orderapi"`
}

func (s *Server) Init(context.Context) error {
	s.Logger().Info("Init")
	r := chi.NewRouter()
	r.Route("/api/orders", func(r chi.Router) {
		r.Post("/", s.CreateOrder)
		r.Get("/{orderid}", s.GetOrderByID)
	})
	//r.Post("/api/orders", s.CreateOrder)
	//r.Get("/api/orders", s.CreateOrder)
	s.handler = r
	return nil
}

// Serve implements the application main.
func Serve(ctx context.Context, s *Server) error {
	s.Logger().Info("OrderAPI listener available.", "addr:", s.orderapi)
	httpServer := &http.Server{
		Handler: s.handler,
	}
	httpServer.Serve(s.orderapi)
	return nil
}
