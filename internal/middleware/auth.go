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

package middleware

import (
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/vpaza/ids/internal/response"
	"github.com/vpaza/ids/pkg/database/models"
)

var log = logger.Logger.WithField("component", "middleware/auth")

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		if sess.Values["cid"] == nil {
			c.Set("x-guest", true)
			c.Set("x-cid", "0")
			c.Set("x-user", nil)
		} else {
			user, err := models.FindUser(sess.Values["cid"])
			if err != nil {
				sess.Values["cid"] = nil
				_ = sess.Save(c.Request(), c.Response())
				c.Set("x-guest", true)
				c.Set("x-cid", "0")
				c.Set("x-user", nil)
			} else {
				c.Set("x-guest", false)
				c.Set("x-cid", sess.Values["cid"])
				c.Set("x-user", user)
			}
		}

		return next(c)
	}
}

func NotGuest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Infof("Got NotGuest middleware=%v", c.Get("x-guest"))
		if c.Get("x-guest").(bool) {
			return response.RespondMessage(c, 401, "Unauthenticated")
		}

		return next(c)
	}
}
