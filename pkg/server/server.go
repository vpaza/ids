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

package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log"
	"golang.org/x/net/http2"

	"github.com/vpaza/ids/internal/binding"
	internalmiddleware "github.com/vpaza/ids/internal/middleware"
	"github.com/vpaza/ids/internal/v1/router"
	"github.com/vpaza/ids/pkg/config"
)

type Server struct {
	E *echo.Echo
}

func NewServer() *Server {
	e := echo.New()

	e.Binder = &binding.CustomBinder{}

	e.HideBanner = true

	generateSecureMiddleware(e)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			return true, nil
		},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	e.Pre(middleware.MethodOverride())
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(session.Middleware(sessions.NewCookieStore([]byte(config.Cfg.Session.Secret))))
	e.Use(internalmiddleware.UpdateCookie)
	e.Use(internalmiddleware.Auth)

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  echoLog.ERROR,
	}))

	return &Server{
		E: e,
	}
}

func (s *Server) handleStart(host string) error {
	switch config.Cfg.Server.Mode {
	case "plain":
		return s.E.Start(host)
	case "tls":
		return s.E.StartTLS(host, config.Cfg.Server.SSLCert, config.Cfg.Server.SSLKey)
	case "h2c":
		sh2 := &http2.Server{
			MaxConcurrentStreams: 250,
			MaxReadFrameSize:     1048576,
			IdleTimeout:          10 * time.Second,
		}
		return s.E.StartH2CServer(host, sh2)
	default:
		return fmt.Errorf("invalid server mode: %s", config.Cfg.Server.Mode)
	}
}

func (s *Server) BuildRoutes() {
	router.SetupRoutes(s.E)

	s.E.Static("/", "static")
}

func (s *Server) Start() {
	go func() {
		if err := s.handleStart(fmt.Sprintf(":%s", config.Cfg.Server.Port)); err != http.ErrServerClosed {
			s.E.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.E.Shutdown(ctx); err != nil {
		s.E.Logger.Fatal(err)
	}
}

func generateSecureMiddleware(e *echo.Echo) {
	c := &middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		ContentSecurityPolicy: "default-src 'self'",
	}

	if config.Cfg.Server.Mode != "plain" {
		c.HSTSExcludeSubdomains = false
		c.HSTSMaxAge = 3600
	}

	e.Use(middleware.SecureWithConfig(*c))
}
