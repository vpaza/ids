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

package sia

import (
	"net/http"
	"time"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/labstack/echo/v4"

	"github.com/vpaza/ids/internal/response"
	"github.com/vpaza/ids/pkg/database/models"
)

type SIADTO struct {
	ATIS             string     `json:"atis"`
	ATISTime         *time.Time `json:"atis_time"`
	DepartureRunways string     `json:"departure_runways"`
	ArrivalRunways   string     `json:"arrival_runways"`
}

func getSIA(e echo.Context) error {
	a, err := models.GetAirport(e.Param("airport"))
	if err != nil {
		return response.RespondError(e, http.StatusInternalServerError, err)
	}

	if a == nil {
		return response.RespondMessage(e, http.StatusNotFound, "Airport Not Found")
	}

	return response.Respond(e, http.StatusOK, a)
}

func patchSIA(e echo.Context) error {
	s := &SIADTO{}
	if err := e.Bind(s); err != nil {
		return response.RespondError(e, http.StatusBadRequest, err)
	}

	a, err := models.GetAirport(e.Param("airport"))
	if err != nil {
		return response.RespondError(e, http.StatusInternalServerError, err)
	}
	if a == nil {
		return response.RespondMessage(e, http.StatusNotFound, "Airport Not Found")
	}

	if s.ATIS != "" {
		a.ATIS = s.ATIS
		n := time.Now()
		a.ATISTime = &n
	}

	if s.DepartureRunways != "" {
		a.DepartureRunways = s.DepartureRunways
	}

	if s.ArrivalRunways != "" {
		a.ArrivalRunways = s.ArrivalRunways
	}

	if err := database.DB.Save(a).Error; err != nil {
		return response.RespondError(e, http.StatusInternalServerError, err)
	}

	return response.RespondBlank(e, http.StatusNotImplemented)
}
