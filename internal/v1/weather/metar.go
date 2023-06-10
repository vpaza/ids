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

package weather

import (
	"fmt"
	"net/http"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/labstack/echo/v4"

	"github.com/vpaza/ids/internal/response"
	"github.com/vpaza/ids/pkg/database/models"
)

type Metar struct {
	Metar  string `json:"metar"`
	MagVar int    `json:"mag_var"`
}

func getMetar(c echo.Context) error {
	var airports []*models.Airport
	if c.Param("airport") == "all" {
		if err := database.DB.Find(&airports).Error; err != nil {
			return response.RespondError(c, http.StatusInternalServerError, err)
		}
		metars := make(map[string]*Metar)
		for _, airport := range airports {
			metars[airport.FAAID] = &Metar{
				Metar:  airport.METAR,
				MagVar: airport.MagVar,
			}
		}

		return response.Respond(c, http.StatusOK, metars)
	}

	airport, err := models.GetAirport(c.Param("airport"))
	if err != nil {
		return response.RespondError(c, http.StatusInternalServerError, err)
	}

	if airport == nil {
		return response.RespondMessage(c, http.StatusNotFound, fmt.Sprintf("Airport not found %v", c.Param("airport")))
	}
	metars := make(map[string]*Metar)
	metars[airport.FAAID] = &Metar{
		Metar:  airport.METAR,
		MagVar: airport.MagVar,
	}
	return response.Respond(c, http.StatusOK, metars)
}
