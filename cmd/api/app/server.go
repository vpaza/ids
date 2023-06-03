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

package app

import (
	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/urfave/cli/v2"

	"github.com/vpaza/ids/internal/middleware"
	"github.com/vpaza/ids/pkg/config"
	"github.com/vpaza/ids/pkg/database/models"
	"github.com/vpaza/ids/pkg/jobs"
	"github.com/vpaza/ids/pkg/oauth"
	"github.com/vpaza/ids/pkg/server"
)

var log = logger.Logger.WithField("component", "server")

func newServerCommand() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start backend server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Value:   "config.yaml",
				Usage:   "Load configuration from `FILE`",
				Aliases: []string{"c"},
				EnvVars: []string{"CONFIG"},
			},
		},
		Action: func(c *cli.Context) error {
			log.Infof("Starting server...")
			log.Infof("config=%s", c.String("config"))

			log.Infof("Loading configuration...")
			_, err := config.ParseConfig(c.String("config"))
			if err != nil {
				return err
			}

			log.Infof("Building web server...")
			srvr := server.NewServer()

			log.Infof("Setting up logger")
			srvr.E.Use(middleware.Logger())

			log.Infof("Building database connection...")
			err = database.Connect(database.DBOptions{
				Host:     config.Cfg.Database.Host,
				Port:     config.Cfg.Database.Port,
				User:     config.Cfg.Database.Username,
				Password: config.Cfg.Database.Password,
				Database: config.Cfg.Database.Database,
				CACert:   config.Cfg.Database.CACert,
				Driver:   "mysql",
				Logger:   logger.Logger,
			})
			if err != nil {
				return err
			}

			log.Infof("Running migrations...")
			err = database.DB.AutoMigrate(
				&models.User{},
				&models.PIREP{},
				&models.Airport{},
				&models.Facility{},
			)
			if err != nil {
				return err
			}

			log.Infof("Building OAuth2 Client...")
			oauth.BuildWithConfig(config.Cfg)

			log.Infof("Building routes...")
			srvr.BuildRoutes()

			log.Infof("Building jobs...")
			jobs.BuildJobs()

			log.Infof("Starting jobs async...")
			jobs.Start()

			log.Infof("Starting web server...")
			srvr.Start()

			return nil
		},
	}
}
