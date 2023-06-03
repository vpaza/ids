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

	"github.com/labstack/echo/v4"

	"github.com/vpaza/ids/internal/response"
	"github.com/vpaza/ids/pkg/database/models"
)

func getMetar(c echo.Context) error {
	airport, err := models.GetAirport(c.Param("airport"))
	if err != nil {
		return response.RespondError(c, http.StatusInternalServerError, err)
	}

	if airport == nil {
		return response.RespondMessage(c, http.StatusNotFound, fmt.Sprintf("Airport not found %v", c.Param("airport")))
	}

	return response.Respond(c, http.StatusOK, airport.METAR)
}
