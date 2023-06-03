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

import (
	"time"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/logger"
	"gorm.io/gorm"
)

var log = logger.Logger.WithField("component", "models/airports")

type Airport struct {
	ID               int64      `json:"id"`
	FAAID            string     `json:"faa_id" gorm:"unique"`
	ICAOID           string     `json:"icao_id" gorm:"unique"`
	ATIS             string     `json:"atis"`
	ATISTime         *time.Time `json:"atis_time"`
	DepartureRunways string     `json:"departure_runways"`
	ArrivalRunways   string     `json:"arrival_runways"`
	METAR            string     `json:"metar"`
	TAF              string     `json:"taf"`
	ParentFacility   int64      `json:"parent_facility"`
}

func GetAirportByICAO(icao string) (*Airport, error) {
	var a Airport

	if err := database.DB.Where("icao_id = ?", icao).First(&a).Error; err != nil {
		return nil, err
	}

	return &a, nil
}

func GetAirportByFAAID(id string) (*Airport, error) {
	var a Airport

	if err := database.DB.Where("faa_id = ?", id).First(&a).Error; err != nil {
		return nil, err
	}

	return &a, nil
}

func GetAirport(id string) (*Airport, error) {
	a, err := GetAirportByICAO(id)
	log.Infof("GetAirportByICAO: %v, %v", a, err)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if a == nil {
		a, err = GetAirportByFAAID(id)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	return a, nil
}
