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

package models

import "time"

type PIREP struct {
	ID           int64      `json:"id"`
	Type         string     `json:"type"`
	Over         string     `json:"over"`
	Time         *time.Time `json:"time"`
	Altitude     int64      `json:"altitude"`
	AircraftType string     `json:"aircraft_type"`
	SkyCover     string     `json:"sky_cover"`
	Weather      string     `json:"weather"`
	Temperature  int64      `json:"temperature"`
	Wind         string     `json:"wind"`
	Turbulence   string     `json:"turbulence"`
	Icing        string     `json:"icing"`
	Remarks      string     `json:"remarks"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
