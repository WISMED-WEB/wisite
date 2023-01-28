package signout

import (
	"fmt"
	"net/http"

	lk "github.com/digisan/logkit"
	so "github.com/digisan/user-mgr/sign-out"
	u "github.com/digisan/user-mgr/user"
	"github.com/labstack/echo/v4"
	"github.com/wismed-web/wisite-api/server/api/sign"
)

// *** after implementing, register with path in 'sign-out.go' *** //

// @Title sign out
// @Summary sign out action.
// @Description
// @Tags    Sign
// @Accept  json
// @Produce json
// @Success 200 "OK - sign-out successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/sign-out/ [get]
// @Security ApiKeyAuth
func SignOut(c echo.Context) error {

	invoker, err := u.Invoker(c)
	if err != nil {
		lk.Warn("%v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	defer invoker.DeleteToken() // only in SignOut calling DeleteToken()

	uname := invoker.UName

	// remove user by 'uname'
	defer sign.UserCache.Delete(uname)

	if err := so.Logout(uname); err != nil {
		lk.Warn("%v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("[%s] sign-out successfully", uname))
}
