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
	"strings"
	"time"

	"github.com/adh-partnership/api/pkg/database"
	"gorm.io/gorm"
)

type Airport struct {
	ID               int64      `json:"id"`
	FAAID            string     `json:"faa_id" gorm:"unique"`
	ICAOID           string     `json:"icao_id" gorm:"unique"`
	ATIS             string     `json:"atis"`
	ATISTime         *time.Time `json:"atis_time"`
	ArrivalATIS      string     `json:"arrival_atis"`
	ArrivalATISTime  *time.Time `json:"arrival_atis_time"`
	DepartureRunways string     `json:"departure_runways"`
	ArrivalRunways   string     `json:"arrival_runways"`
	METAR            string     `json:"metar"`
	TAF              string     `json:"taf"`
	MagVar           int        `json:"mag_var"`
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

func GetAirports(list string) ([]*Airport, error) {
	var airports []*Airport

	if list == "all" || list == "" {
		if err := database.DB.Find(&airports).Error; err != nil {
			return nil, err
		}

		return airports, nil
	}

	if err := database.DB.Where("icao_id IN ?", strings.Split(list, ",")).Or("faa_id IN ?", strings.Split(list, ",")).Find(&airports).Error; err != nil {
		return nil, err
	}

	return airports, nil
}
