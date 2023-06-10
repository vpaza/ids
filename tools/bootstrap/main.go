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

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/logger"

	"github.com/vpaza/ids/pkg/config"
	"github.com/vpaza/ids/pkg/database/models"
	"github.com/vpaza/ids/pkg/utils"
)

type Airport struct {
	Name string `json:"name"`
}

type Facility struct {
	Airports   []*Airport `json:"airports"`
	WXStations []string   `json:"wx_stations"`
}

func main() {
	cfg, err := config.ParseConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	err = database.Connect(database.DBOptions{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.Username,
		Password: cfg.Database.Password,
		Database: cfg.Database.Database,
		CACert:   cfg.Database.CACert,
		Driver:   "mysql",
		Logger:   logger.Logger,
	})
	if err != nil {
		// Wait 5 seconds and try again, and keep doing it
		// until we can connect
		for {
			fmt.Println("Error connecting to database, retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			err = database.Connect(database.DBOptions{
				Host:     cfg.Database.Host,
				Port:     cfg.Database.Port,
				User:     cfg.Database.Username,
				Password: cfg.Database.Password,
				Database: cfg.Database.Database,
				CACert:   cfg.Database.CACert,
				Driver:   "mysql",
				Logger:   logger.Logger,
			})
			if err == nil {
				break
			}
		}
	}

	err = database.DB.AutoMigrate(
		&models.User{},
		&models.PIREP{},
		&models.Airport{},
	)
	if err != nil {
		panic(err)
	}

	facility := &Facility{}
	f, err := os.Open("frontend/src/facility.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &facility)
	if err != nil {
		panic(err)
	}

	// Truncate the airport table
	err = database.DB.Exec("TRUNCATE TABLE airports").Error
	if err != nil {
		panic(err)
	}

	for _, airport := range facility.Airports {
		for {
			fmt.Printf("Handling airport %s... ", airport.Name)
			ret, err := lookupAirportInfo(airport.Name)
			if err != nil {
				fmt.Printf("Error lookuping up icao identifier for %s: %s. Sleeping and retrying...\n", airport.Name, err)
				time.Sleep(5 * time.Second)
				continue
			}
			a := &models.Airport{
				FAAID:            airport.Name,
				ICAOID:           ret.ICAOIdent,
				ATIS:             "",
				ATISTime:         utils.Now(),
				DepartureRunways: "",
				ArrivalRunways:   "",
				METAR:            "",
				TAF:              "",
				MagVar:           ret.MagVar,
			}

			// Recreate
			err = database.DB.Create(a).Error
			if err != nil {
				fmt.Printf("Error creating airport %s: %s. Skipping.\n", airport.Name, err)
				break
			}

			fmt.Printf("done\n")
			break
		}
		time.Sleep(250 * time.Millisecond)
	}

	for _, airport := range facility.WXStations {
		for {
			fmt.Printf("Handling airport %s... ", airport)
			// If it exists, skip
			var count int64
			err := database.DB.Model(&models.Airport{}).Where("icao_id = ?", airport).Count(&count).Error
			if err != nil {
				fmt.Printf("Error checking if airport %s exists: %s. Skipping.\n", airport, err)
				panic(err)
			}
			if count > 0 {
				fmt.Printf("already exists.. skipping\n")
				break
			}

			ret, err := lookupAirportInfo(airport)
			if err != nil {
				fmt.Printf("Error lookuping up icao identifier for %s: %s. Sleeping and retrying...\n", airport, err)
				time.Sleep(5 * time.Second)
				continue
			}

			// Recreate
			err = database.DB.Create(&models.Airport{
				FAAID:  ret.FAAIdent,
				ICAOID: airport,
				MagVar: ret.MagVar,
			}).Error
			if err != nil {
				fmt.Printf("Error creating airport %s: %s. Skipping.\n", airport, err)
				continue
			}
			fmt.Printf("done\n")
			break
		}
		time.Sleep(250 * time.Millisecond)
	}

	fmt.Println("Done!")
}

// We only need certain fields...
type AviationAPIResponse struct {
	FAAIdent          string `json:"faa_ident"`
	ICAOIdent         string `json:"icao_ident"`
	MagneticVariation string `json:"magnetic_variation"`
}

type Ret struct {
	FAAIdent  string
	ICAOIdent string
	MagVar    int
}

func lookupAirportInfo(faaid string) (*Ret, error) {
	// Query https://api.aviation.api.com/v1/airports?apt=
	// and return the ICAO code
	r, err := http.NewRequest("GET", fmt.Sprintf("https://api.aviationapi.com/v1/airports?apt=%s", faaid), nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp map[string][]*AviationAPIResponse
	err = json.Unmarshal(contents, &apiResp)
	if err != nil {
		fmt.Printf("Unable to unmarshal response: %s", string(contents))
		return nil, err
	}

	if len(apiResp) > 0 {
		return &Ret{
			FAAIdent:  apiResp[faaid][0].FAAIdent,
			ICAOIdent: apiResp[faaid][0].ICAOIdent,
			MagVar:    calcMagVar(apiResp[faaid][0].MagneticVariation),
		}, nil
	}

	return nil, fmt.Errorf("no icao code found for %s", faaid)
}

func calcMagVar(mv string) int {
	// Convert string of 15E to -15 or 15W to 15
	var m int
	if mv[len(mv)-1:] == "E" {
		m, _ = strconv.Atoi(mv[:len(mv)-1])
		m *= -1
	} else {
		m, _ = strconv.Atoi(mv[:len(mv)-1])
	}

	return m
}
