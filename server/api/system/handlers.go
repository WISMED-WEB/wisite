package system

import (
	"net/http"

	"github.com/digisan/gotk/project"
	"github.com/labstack/echo/v4"
)

// *** after implementing, register with path in 'system.go' *** //

// @Title api service version
// @Summary get this api service version
// @Description
// @Tags    system
// @Accept  json
// @Produce json
// @Success 200 "OK - get its version"
// @Failure 500 "Fail - default version 'v0.0.0' applies"
// @Router /api/system/ver [get]
func Ver(c echo.Context) error {
	ver, ok := project.GitVer("v0.0.0")
	if ok {
		return c.String(http.StatusOK, ver)
	}
	return c.String(http.StatusInternalServerError, "failed to get version")
}

// @Title api service tag
// @Summary get this api service project github version tag
// @Description
// @Tags    system
// @Accept  json
// @Produce json
// @Success 200 "OK - get its tag"
// @Failure 500 "Fail - couldn't get service project tag"
// @Router /api/system/ver-tag [get]
func VerTag(c echo.Context) error {
	tag, err := project.GitTag()
	if err == nil {
		return c.String(http.StatusOK, tag)
	}
	return c.String(http.StatusInternalServerError, "failed to get version tag")
}
