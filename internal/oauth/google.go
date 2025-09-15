package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"pets_rest/pkg/helper"

	"github.com/gofiber/fiber/v3/middleware/session"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	stateKey = "oauth_state"
	pkceKey  = "oauth_pkce"
)

type GoogleProvider struct {
	cfg *oauth2.Config
}

func NewGoogle() *GoogleProvider {
	return &GoogleProvider{
		cfg: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),     // to be set from config
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"), // to be set from config
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),  // to be set from config
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (p *GoogleProvider) AuthURLWithPKCEandState(sess *session.Session) (string, error) {
	state := NewState()
	ver, ch := NewPKCE()

	sess.Set(stateKey, state)
	sess.Set(pkceKey, ver)

	if err := sess.Save(); err != nil {
		return "", err
	}

	return p.cfg.AuthCodeURL(state,
		oauth2.SetAuthURLParam("code_challenge", ch),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	), nil
}

func (p *GoogleProvider) HandleCallback(ctx context.Context, sess *session.Session, state, code string) (User, error) {
	if state == "" || state != helper.GetString(sess.Get(stateKey)) {
		return User{}, errors.New("invalid state parameter")
	}

	ver := helper.GetString(sess.Get(pkceKey))
	if ver == "" {
		return User{}, errors.New("invalid PKCE verifier")
	}

	tok, err := p.cfg.Exchange(ctx, code,
		oauth2.SetAuthURLParam("code_verifier", ver),
	)
	if err != nil {
		return User{}, err
	}

	client := p.cfg.Client(ctx, tok)
	resp, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	var ui struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ui); err != nil {
		return User{}, err
	}

	return User{
		Provider:      "google",
		ProviderID:    ui.Sub,
		Email:         ui.Email,
		EmailVerified: ui.EmailVerified,
		Name:          ui.Name,
		AvatarURL:     ui.Picture,
	}, nil
}
