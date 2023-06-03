/*
 * Copyright Daniel Hawton
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

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
		RedirectURL:  fmt.Sprintf("%s%s", c.OAuth.MyBaseURL, "/v1/oauth/callback"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s%s", c.OAuth.BaseURL, config.Cfg.OAuth.EndpointAuthorize),
			TokenURL: fmt.Sprintf("%s%s", c.OAuth.BaseURL, config.Cfg.OAuth.EndpointToken),
		},
	}
}
