package oauth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anfastk/mergespace/auth/internal/auth/application/dto"

	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

type GitHubOAuthProvider struct {
	config *oauth2.Config
}

func NewGitHubOAuthProvider(clientID string, clientSecret string, redirectURL string) *GitHubOAuthProvider {

	return &GitHubOAuthProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"read:user",
				"user:email",
			},
			Endpoint: githuboauth.Endpoint,
		},
	}
}

func (g *GitHubOAuthProvider) GetLoginURL(state string) string {

	return g.config.AuthCodeURL(state)
}

func (g *GitHubOAuthProvider) GetGitHubUser(ctx context.Context, code string) (*dto.GitHubUser, error) {

	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := g.config.Client(ctx, token)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user dto.GitHubUser

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	if user.Email == "" {

		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err != nil {
			return nil, err
		}
		defer emailResp.Body.Close()

		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}

		if err := json.NewDecoder(emailResp.Body).Decode(&emails); err != nil {
			return nil, err
		}

		for _, e := range emails {
			if e.Primary && e.Verified {
				user.Email = e.Email
				break
			}
		}
	}

	if user.Email == "" {
		return nil, fmt.Errorf("github email not available")
	}

	return &user, nil
}
