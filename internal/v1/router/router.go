package router

import "github.com/labstack/echo/v4"

var routeGroups map[string]func(e *echo.Echo)

