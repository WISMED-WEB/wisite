package signout

import (
	"fmt"
	"net/http"

	so "github.com/digisan/user-mgr/sign-out"
	"github.com/labstack/echo/v4"
)

// *** after implementing, register with path in 'api_reg.go' ***

// SignIn godoc
// @Title sign out
// @Summary sign out action.
// @Description
// @Tags    sign-out
// @Accept  json
// @Produce json
// @Param   uname query string true "unique user name"
// @Success 200 "OK - sign-out successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/sign-out/signout [get]
func SignOut(c echo.Context) error {
	if err := so.Logout(c.QueryParam("uname")); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln(err))
	}
	return c.String(http.StatusOK, "sign-out successfully")
}
