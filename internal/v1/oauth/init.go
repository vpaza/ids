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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/vpaza/ids/internal/response"
	"github.com/vpaza/ids/pkg/config"
	"github.com/vpaza/ids/pkg/database/models"
	"github.com/vpaza/ids/pkg/oauth"
)

var log = logger.Logger.WithField("component", "oauth")

type SSOUserInfoResponse struct {
	Message string `json:"message"`
	User    struct {
		CID       uint   `json:"cid"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	} `json:"user"`
}

func SetupRoutes(g *echo.Group) {
	g.GET("/callback", getCallback)
	g.GET("/login", getLogin)
	g.GET("/logout", getLogout)
}

func getCallback(e echo.Context) error {
	sess, _ := session.Get("session", e)
	if sess.Values["state"] != e.QueryParam("state") {
		log.Warnf("State is not equal: %s != %s", sess.Values["state"], e.QueryParam("state"))
		return response.RespondMessage(e, http.StatusForbidden, "Forbidden")
	}

	token, err := oauth.OAuthConfig.Exchange(e.Request().Context(), e.QueryParam("code"))
	if err != nil {
		log.Errorf("Error exchanging token: %v", err)
		return response.RespondMessage(e, http.StatusInternalServerError, "Internal Server Error")
	}
	res, err := http.NewRequest("GET", fmt.Sprintf("%s%s", config.Cfg.OAuth.BaseURL, config.Cfg.OAuth.EndpointUserInfo), nil)
	res.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	res.Header.Add("Accept", "application/json")
	res.Header.Add("User-Agent", "adh-partnership-api")
	if err != nil {
		return response.RespondMessage(e, http.StatusInternalServerError, "Internal Server Error")
	}

	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		return response.RespondMessage(e, http.StatusInternalServerError, "Internal Server Error")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.RespondMessage(e, http.StatusInternalServerError, "Internal Server Error")
	}

	if resp.StatusCode >= 299 {
		log.Warnf("Error getting user info: %s, %s", resp.Status, string(contents))
		return response.RespondMessage(e, http.StatusForbidden, "Forbidden")

	}

	user := &SSOUserInfoResponse{}
	if err := json.Unmarshal(contents, &user); err != nil {
		return response.RespondMessage(e, http.StatusInternalServerError, "Internal Server Error")
	}

	u, err := models.FindUser(user.User.CID)
	if err != nil {
		return response.RespondMessage(e, http.StatusInternalServerError, "Internal Server Error")
	}

	if u == nil {
		u = &models.User{
			CID:       user.User.CID,
			FirstName: user.User.FirstName,
			LastName:  user.User.LastName,
			Email:     user.User.Email,
		}
		if err := database.DB.Create(&u); err != nil {
			return response.RespondMessage(e, http.StatusInternalServerError, "Internal Server Error")
		}
	}

	sess.Values["cid"] = u.CID
	sess.Values["state"] = nil
	_ = sess.Save(e.Request(), e.Response())

	redirect := sess.Values["redirect"]
	if redirect != nil {
		return e.Redirect(http.StatusTemporaryRedirect, redirect.(string))
	}

	return response.RespondMessage(e, http.StatusOK, "Logged In")
}

func getLogin(e echo.Context) error {
	state, _ := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 64)
	sess, _ := session.Get("session", e)
	sess.Values["state"] = state
	_ = sess.Save(e.Request(), e.Response())

	return e.Redirect(
		http.StatusTemporaryRedirect,
		oauth.OAuthConfig.AuthCodeURL(state),
	)
}

func getLogout(e echo.Context) error {
	sess, _ := session.Get("session", e)
	sess.Values["cid"] = nil
	_ = sess.Save(e.Request(), e.Response())

	return response.RespondBlank(e, http.StatusNoContent)
}
