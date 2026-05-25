package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	httpOAuth "github.com/anfastk/mergespace/auth/internal/auth/adapter/inbound/http/oauth"
	"github.com/anfastk/mergespace/auth/internal/auth/infrastructure/di"
	"github.com/anfastk/mergespace/contracts/gen/go/proto/auth/v1/authv1connect"
)

func main() {

	app := di.BuildApp()

	go app.Worker.Start(context.Background())

	r := chi.NewRouter()

	_, handler := authv1connect.NewAuthServiceHandler(app.Handler)

	r.Mount("/", handler)

	googleHandler := httpOAuth.NewGoogleHandler(
		app.HandlerUsecase,
		app.GoogleProvider,
	)

	githubHandler := httpOAuth.NewGitHubHandler(
		app.HandlerUsecase,
		app.GitHubProvider,
	)

	r.Get("/auth/github/login", githubHandler.Login)
	r.Get("/auth/github/callback", githubHandler.Callback)

	r.Get("/auth/google/login", googleHandler.Login)
	r.Get("/auth/google/callback", googleHandler.Callback)

	log.Println("Auth Service running on :8080")

	if err := http.ListenAndServe(
		":8080",
		r,
	); err != nil {
		log.Fatal(err)
	}
}
