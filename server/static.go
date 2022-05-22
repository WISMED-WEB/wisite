package main

import "github.com/labstack/echo/v4"

func hookStatic(e *echo.Echo) {
	// e.Static("/", "www")          // host www folder, allow js/ etc.
	// e.File("/", "www/index.html") // host www/index.html static file

	e.Static("/", "data/user-space")
}
