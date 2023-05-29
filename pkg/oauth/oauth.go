package oauth

import (
	"fmt"

	"golang.org/x/oauth2"

	"github.com/vpaza/ids/pkg/config"
)

var OAuthConfig *oauth2.Config

func BuildWithConfig(c *config.Config) {
	OAuthConfig = &oauth2.Config{
		ClientID:     c.OAuth.ClientID,
		ClientSecret: c.OAuth.ClientSecret,
		Scopes:       []string{"identify", "email"},
		RedirectURL:  fmt.Sprintf("%s%s", c.OAuth.MyBaseURL, "/v1/user/login/callback"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s%s", c.OAuth.BaseURL, config.Cfg.OAuth.EndpointAuthorize),
			TokenURL: fmt.Sprintf("%s%s", c.OAuth.BaseURL, config.Cfg.OAuth.EndpointToken),
		},
	}
}
