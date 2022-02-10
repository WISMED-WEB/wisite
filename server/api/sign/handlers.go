package sign

import (
	"fmt"
	"net/http"
	"sync"

	lk "github.com/digisan/logkit"
	si "github.com/digisan/user-mgr/sign-in"
	so "github.com/digisan/user-mgr/sign-out"
	su "github.com/digisan/user-mgr/sign-up"
	usr "github.com/digisan/user-mgr/user"
	"github.com/labstack/echo/v4"
	md "github.com/wismed-web/wisite/module"
)

// *** after implementing, register with path in 'api_sign.go' ***

// @Title sign in
// @Summary sign in action. if ok, got token
// @Description
// @Tags    sign
// @Accept  json
// @Produce json
// @Param   uname query string true "unique user name"
// @Param   pwd   query string true "password"
// @Success 200 "OK - sign-in successfully"
// @Failure 400 "Fail - incorrect password"
// @Failure 500 "Fail - internal error"
// @Router /api/sign/in [get]
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

// @Title sign out
// @Summary sign out action.
// @Description
// @Tags    sign
// @Accept  json
// @Produce json
// @Param   uname query string true "unique user name"
// @Success 200 "OK - sign-out successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/sign/out [get]
func SignOut(c echo.Context) error {
	if err := so.Logout(c.QueryParam("uname")); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "sign-out successfully")
}

var (
	mUser = &sync.Map{} // users waiting for verifying email code
)

// @Title register a new user
// @Summary send user's basic info for registry
// @Description
// @Tags    sign
// @Accept  multipart/form-data
// @Produce json
// @Param   uname   formData   string  true  "unique user name"
// @Param   email   formData   string  true  "user's email"
// @Param   name    formData   string  true  "user's real name"
// @Param   pwd     formData   string  true  "user's password"
// @Success 200 "OK - then waiting for verification code"
// @Failure 400 "Fail - invalid registry fields"
// @Failure 500 "Fail - internal error"
// @Router /api/sign/new [post]
func NewUser(c echo.Context) error {

	// lk.Debug("[%v] [%v] [%v] [%v]", c.FormValue("uname"), c.FormValue("email"), c.FormValue("name"), c.FormValue("pwd"))

	user := usr.User{
		Active:   "T",
		UName:    c.FormValue("uname"),
		Email:    c.FormValue("email"),
		Name:     c.FormValue("name"),
		Password: c.FormValue("pwd"),
		Regtime:  "TBD",
		MemLevel: "0",
	}

	// su.SetValidator(map[string]func(string) bool{ })

	lk.Log("%v", user)

	if err := su.ChkInput(user); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprint(err))
	}
	if err := su.ChkEmail(user); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprint(err))
	}

	mUser.Store(user.UName, user)

	return c.String(http.StatusOK, "waiting verification code in your email")
}

// @Title verify new user's email
// @Summary send back email verification code
// @Description
// @Tags    sign
// @Accept  multipart/form-data
// @Produce json
// @Param   uname  formData  string  true  "unique user name"
// @Param   code   formData  string  true  "verification code (in user's email)"
// @Success 200 "OK - sign-up successfully"
// @Failure 400 "Fail - incorrect verification code"
// @Failure 500 "Fail - internal error"
// @Router /api/sign/verify-email [post]
func VerifyEmail(c echo.Context) error {

	code := c.FormValue("code")
	uname := c.FormValue("uname")

	user, ok := mUser.LoadAndDelete(uname)
	if !ok {
		return c.String(http.StatusBadRequest, "need re-sending verification code")
	}

	if err := su.VerifyCode(user.(usr.User), code); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprint(err))
	}

	// store into db
	if err := su.Store(user.(usr.User)); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprint(err))
	}

	// sign-up ok calling...
	{
		md.AllocDisk(uname)
	}

	return c.String(http.StatusOK, "registered successfully")
}
