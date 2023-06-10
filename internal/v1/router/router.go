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

package router

import (
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/labstack/echo/v4"

	"github.com/vpaza/ids/internal/v1/external"
	"github.com/vpaza/ids/internal/v1/oauth"
	"github.com/vpaza/ids/internal/v1/sia"
	"github.com/vpaza/ids/internal/v1/weather"
)

var (
	log         = logger.Logger.WithField("component", "router")
	routeGroups map[string]func(e *echo.Group)
)

func init() {
	routeGroups = make(map[string]func(e *echo.Group))
	routeGroups["/external"] = external.Routes
	routeGroups["/oauth"] = oauth.Routes
	routeGroups["/weather"] = weather.Routes
	routeGroups["/sia"] = sia.Routes
}

func SetupRoutes(e *echo.Echo) {
	e.GET("/health", healthCheckHandler)
	e.GET("/ready", readyCheckHandler)

	v1 := e.Group("/v1")
	for prefix, group := range routeGroups {
		log.Infof("Loading route prefix: %s", prefix)
		grp := v1.Group(prefix)
		group(grp)
	}
}
