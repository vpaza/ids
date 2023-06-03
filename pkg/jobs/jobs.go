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

package jobs

import (
	"time"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/go-co-op/gocron"

	"github.com/vpaza/ids/pkg/database/models"
	"github.com/vpaza/ids/pkg/weather"
)

var s *gocron.Scheduler

var log = logger.Logger.WithField("component", "jobs")

func BuildJobs() {
	s = gocron.NewScheduler(time.UTC)
	_, _ = s.Every(5).Minutes().SingletonMode().Do(updateWeather)
	_, _ = s.Every(1).Minutes().SingletonMode().Do(cleanPireps)
}

func Start() {
	s.StartAsync()
}

func updateWeather() {
	log.Debugf("Running: updateWeather")

	var airports []*models.Airport

	if err := database.DB.Find(&airports).Error; err != nil {
		log.Errorf("Error getting airports: %s", err)
		return
	}

	for _, airport := range airports {
		metar, err := weather.GetMetar(airport.ICAOID)
		if err != nil {
			log.Errorf("Error getting METAR for %s: %s", airport.ICAOID, err)
			continue
		}
		airport.METAR = metar.RawText
		if err := database.DB.Save(airport).Error; err != nil {
			log.Errorf("Error saving airport: %s", err)
			continue
		}
	}
}

func cleanPireps() {
	log.Debugf("Running: cleanPireps")

	// Find and delete any pirep older than 2 hours
	if err := database.DB.Where("created_at < ?", time.Now().Add(-2*time.Hour)).Delete(&models.PIREP{}).Error; err != nil {
		log.Errorf("Error cleaning PIREPs: %s", err)
	}
}
