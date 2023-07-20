package main

import (
	"context"
	"fmt"

	"github.com/ServiceWeaver/weaver"
)

// Reverser component.
type Greeter interface {
	Greet(context.Context, string) (string, error)
}

// Implementation of the Reverser component.
type greeter struct {
	weaver.Implements[Greeter]
}

func (r *greeter) Greet(_ context.Context, name string) (string, error) {
	str := fmt.Sprintf("Hello, %s, welcome to Service Weaver!", name)
	return str, nil
}
