package main

import (
	"context"
	"log"

	"github.com/ServiceWeaver/weaver"

	"service-weaver/orderapp/orderservice"
)

func main() {
	if err := weaver.Run(context.Background(), orderservice.Serve); err != nil {
		log.Fatal(err)
	}
}
