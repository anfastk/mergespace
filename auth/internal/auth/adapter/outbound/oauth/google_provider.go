package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/anfastk/mergespace/auth/internal/auth/application/dto"
)

type GoogleOAuthProvider struct {
	config *oauth2.Config
}

func NewGoogleOAuthProvider(clientID string, clientSecret string, redirectURL string) *GoogleOAuthProvider {

	return &GoogleOAuthProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"openid",
				"profile",
				"email",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (g *GoogleOAuthProvider) GetLoginURL(state string) string {
	return g.config.AuthCodeURL(state)
}

func (g *GoogleOAuthProvider) GetGoogleUser(ctx context.Context, code string) (*dto.GoogleUser, error) {

	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := g.config.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch google user")
	}

	var user dto.GoogleUser

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
