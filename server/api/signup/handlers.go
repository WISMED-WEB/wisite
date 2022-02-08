package module1

import (
	"fmt"
	"net/http"

	lk "github.com/digisan/logkit"
	su "github.com/digisan/user-mgr/sign-up"
	usr "github.com/digisan/user-mgr/user"
	"github.com/labstack/echo/v4"
)

// after implementing, register with path in 'api_reg.go'

// NewUser godoc
// @Title register a new user
// @Summary send user's basic info for registry
// @Description
// @Tags    sign-up
// @Accept  multipart/form-data
// @Produce json
// @Param   uname   formData   string  true  "unique user name"
// @Param   email   formData   string  true  "user's email"
// @Param   name    formData   string  true  "user's real name"
// @Param   pwd     formData   string  true  "user's password"
// @Success 200 "OK - then waiting for verification code"
// @Failure 400 "Fail - invalid registry fields"
// @Failure 500 "Fail - internal error"
// @Router /api/sign-up/new [post]
func NewUser(c echo.Context) error {

	lk.Debug("[%v] [%v] [%v] [%v]", c.FormValue("uname"), c.FormValue("email"), c.FormValue("name"), c.FormValue("pwd"))

	user := usr.User{
		Active:   "T",
		UName:    c.FormValue("uname"),
		Email:    c.FormValue("email"),
		Name:     c.FormValue("name"),
		Password: c.FormValue("pwd"),
		Regtime:  "TBD",
	}

	lk.Log("%v", user)

	if err := su.ChkInput(user); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprint(err))
	}
	if err := su.ChkEmail(user); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprint(err))
	}
	return c.String(http.StatusOK, "waiting verification code in your email")
}

// -- @Title new user's email verification
// -- @Description verify user's email by checking code
// -- @Success 200 "OK"
// -- @Failure 400 "Fail"
// -- @Failure 500 "Fail"
// -- @Router /api/sign-up/verify-email [post]
// func VerifyEmail(c echo.Context) error {

// 	code := c.Get("code").(string)
// 	user := usr.User{UName: c.Get("uname").(string)}

// 	if err := su.VerifyCode(user, code); err != nil {
// 		fmt.Println("Sign-Up failed:", err)
// 		return
// 	}

// 	// store into db
// 	if err := su.Store(user); err != nil {
// 		return c.String(http.StatusInternalServerError, fmt.Sprint(err))
// 	}
// 	return c.String(http.StatusOK, "registered successfully")
// }
