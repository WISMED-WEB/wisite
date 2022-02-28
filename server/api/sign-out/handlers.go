package signout

import (
	"fmt"
	"net/http"

	so "github.com/digisan/user-mgr/sign-out"
	usr "github.com/digisan/user-mgr/user"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/wismed-web/wisite-api/server/api/sign"
)

// *** after implementing, register with path in 'sign-out.go' *** //

// @Title sign out
// @Summary sign out action.
// @Description
// @Tags    sign
// @Accept  json
// @Produce json
// @Success 200 "OK - sign-out successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/sign-out/ [get]
func SignOut(c echo.Context) error {
	userTkn := c.Get("user").(*jwt.Token)
	claims := userTkn.Claims.(*usr.UserClaims)
	defer claims.DeleteToken()

	uname := claims.UName
	if err := so.Logout(uname); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// remove user space for 'uname'
	sign.MapUserSpace.Delete(uname)

	return c.String(http.StatusOK, fmt.Sprintf("[%s] sign-out successfully", uname))
}