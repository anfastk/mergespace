package main

import (
	"log"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/anfastk/mergespace/auth/internal/auth/infrastructure/di"
	"github.com/anfastk/mergespace/contracts/gen/go/proto/auth/v1/authv1connect"
)

func main() {
	handler := di.BuildAuthHandler()

	mux := http.NewServeMux()

	path, h := authv1connect.NewAuthServiceHandler(handler)
	log.Println("Registered Connect path:", path)
	mux.Handle(path, h)

	log.Println("Auth Service (Connect) running on :8080")

	if err := http.ListenAndServe(
		":8080",
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		log.Fatal(err)
	}
}
