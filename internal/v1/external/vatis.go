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

package external

import (
	"fmt"
	"io"
	"net/http"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/labstack/echo/v4"

	"github.com/vpaza/ids/internal/response"
	"github.com/vpaza/ids/pkg/database/models"
	"github.com/vpaza/ids/pkg/utils"
)

type vATISDTO struct {
	Facility          string `json:"facility"`
	Preset            string `json:"preset"`
	ATISLetter        string `json:"atis_letter"`
	ATISType          string `json:"atis_type"`
	AirportConditions string `json:"airport_conditions"`
	NOTAMs            string `json:"notams"`
	Timestamp         string `json:"timestamp"`
	Version           string `json:"version"`
}

var log = logger.Logger.WithField("component", "vatis")

func postvATIS(c echo.Context) error {
	// Dump the body
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return response.RespondError(c, http.StatusBadRequest, err)
	}
	log.Infof("Received vATIS: %v", string(body))

	var vatis vATISDTO
	if err := c.Bind(&vatis); err != nil {
		return response.RespondError(c, http.StatusBadRequest, err)
	}

	airport, err := models.GetAirport(vatis.Facility)
	if err != nil {
		return response.RespondError(c, http.StatusInternalServerError, err)
	}

	if airport == nil {
		return response.RespondMessage(c, http.StatusNotFound, fmt.Sprintf("Airport not found %v", c.Param("airport")))
	}

	log.Infof("Received vATIS for %v, %+v", airport.FAAID, vatis)

	if vatis.ATISType == "arrival" {
		airport.ArrivalATIS = vatis.ATISLetter
		airport.ArrivalATISTime = utils.Now()
	} else {
		airport.ATIS = vatis.ATISLetter
		airport.ATISTime = utils.Now()
	}

	log.Infof("Setting airport data: %+v", airport)

	if err := database.DB.Save(airport).Error; err != nil {
		return response.RespondError(c, http.StatusInternalServerError, err)
	}

	return response.Respond(c, http.StatusOK, airport.METAR)
}
