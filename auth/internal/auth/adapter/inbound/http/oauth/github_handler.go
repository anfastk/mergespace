package oauth

import (
	"net/http"

	"github.com/anfastk/mergespace/auth/internal/auth/adapter/outbound/oauth"
	"github.com/anfastk/mergespace/auth/internal/auth/application/port/inbound"
)

type GitHubHandler struct {
	authService inbound.AuthUseCase
	provider    *oauth.GitHubOAuthProvider
}

func NewGitHubHandler(
	authService inbound.AuthUseCase, provider *oauth.GitHubOAuthProvider) *GitHubHandler {

	return &GitHubHandler{
		authService: authService,
		provider:    provider,
	}
}

func (h *GitHubHandler) Login(
	w http.ResponseWriter,
	r *http.Request,
) {

	url := h.provider.GetLoginURL("random-state")

	http.Redirect(
		w,
		r,
		url,
		http.StatusTemporaryRedirect,
	)
}

func (h *GitHubHandler) Callback(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")

	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	res, err := h.authService.GitHubLogin(
		r.Context(),
		code,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(res.AccessToken))
}
