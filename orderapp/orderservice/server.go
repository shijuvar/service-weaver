package orderservice

import (
	"context"
	"net/http"

	"github.com/ServiceWeaver/weaver"
	chi "github.com/go-chi/chi/v5"

	"github.com/shijuvar/service-weaver/orderapp/notificationservice"
	"github.com/shijuvar/service-weaver/orderapp/paymentservice"
)

type Server struct {
	weaver.Implements[weaver.Main]

	handler http.Handler // http router instance

	paymentService      weaver.Ref[paymentservice.Service]
	notificationService weaver.Ref[notificationservice.Service]
	//orderRepository     weaver.Ref[cockroachdb.Repository]

	orderapi weaver.Listener //`weaver:"orderapi"`
}

func (s *Server) Init(context.Context) error {
	s.Logger().Info("Init")
	r := chi.NewRouter()
	r.Route("/api/orders", func(r chi.Router) {
		r.Post("/", s.CreateOrder)
		r.Get("/{orderid}", s.GetOrderByID)
	})
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
