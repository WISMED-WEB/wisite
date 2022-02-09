package signin

import (
	"fmt"
	"net/http"

	lk "github.com/digisan/logkit"
	si "github.com/digisan/user-mgr/sign-in"
	usr "github.com/digisan/user-mgr/user"
	"github.com/labstack/echo/v4"
)

// *** after implementing, register with path in 'api_reg.go' ***

// SignIn godoc
// @Title sign in
// @Summary sign in action. if ok, got token
// @Description
// @Tags    sign-in
// @Accept  json
// @Produce json
// @Param   uname query string true "unique user name"
// @Param   pwd   query string true "password"
// @Success 200 "OK - sign-in successfully"
// @Failure 400 "Fail - incorrect password"
// @Failure 500 "Fail - internal error"
// @Router /api/sign-in/signin [get]
func SignIn(c echo.Context) error {

	lk.Debug("[%v] [%v]", c.QueryParam("uname"), c.QueryParam("pwd"))

	user := usr.User{
		UName:    c.QueryParam("uname"),
		Password: c.QueryParam("pwd"),
	}

	if err := si.UserExists(user); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprint(err))
	}

	if !si.PwdOK(user) {
		return c.String(http.StatusBadRequest, "password incorrect")
	}

	defer lk.FailOnErr("%v", si.Trail(user.UName))

	claims := usr.MakeUserClaims(user)
	token := claims.GenToken()
	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
