package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ServiceWeaver/weaver"
)

func main() {
	if err := weaver.Run(context.Background(), serve); err != nil {
		log.Fatal(err)
	}
}

type app struct {
	weaver.Implements[weaver.Main]
	reverser weaver.Ref[Reverser]
	hello    weaver.Listener `weaver:"hello"`
}

func serve(ctx context.Context, app *app) error {
	// The hello listener will listen on a random port chosen by the operating
	// system. This behavior can be changed in the config file.
	fmt.Printf("hello listener available on %v\n", app.hello)
	// Serve the /hello endpoint.
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "World"
		}
		reversed, err := app.reverser.Get().Reverse(ctx, name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Hello, %s!\n", reversed)
	})
	return http.Serve(app.hello, nil)
}
