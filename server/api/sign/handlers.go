package sign

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	fm "github.com/digisan/file-mgr"
	lk "github.com/digisan/logkit"
	rp "github.com/digisan/user-mgr/reset-pwd"
	si "github.com/digisan/user-mgr/sign-in"
	su "github.com/digisan/user-mgr/sign-up"
	u "github.com/digisan/user-mgr/user"
	"github.com/labstack/echo/v4"
)

// *** after implementing, register with path in 'sign.go' *** //

var (
	UserCache = &sync.Map{} // map[string]*u.User, *** record logged-in user ***
)

// @Title register a new user
// @Summary sign up action, step 1. send user's basic info for registry
// @Description
// @Tags    Sign
// @Accept  multipart/form-data
// @Produce json
// @Param   uname   formData   string  true  "unique user name"
// @Param   email   formData   string  true  "user's email" Format(email)
// @Param   name    formData   string  true  "user's real full name"
// @Param   pwd     formData   string  true  "user's password"
// @Success 200 "OK - then waiting for verification code"
// @Failure 400 "Fail - invalid registry fields"
// @Failure 500 "Fail - internal error"
// @Router /api/sign/new [post]
func NewUser(c echo.Context) error {

	// lk.Debug("[%v] [%v] [%v] [%v]", c.FormValue("uname"), c.FormValue("email"), c.FormValue("name"), c.FormValue("pwd"))

	user := &u.User{
		Core: u.Core{
			UName:    c.FormValue("uname"),
			Email:    c.FormValue("email"),
			Password: c.FormValue("pwd"),
		},
		Profile: u.Profile{
			Name:           c.FormValue("name"),
			Phone:          "",
			Country:        "",
			City:           "",
			Addr:           "",
			PersonalIDType: "",
			PersonalID:     "",
			Gender:         "",
			DOB:            "",
			Position:       "",
			Title:          "",
			Employer:       "",
			Bio:            "",
			AvatarType:     "",
			Avatar:         []byte{},
		},
		Admin: u.Admin{
			RegTime:   time.Now().Truncate(time.Second),
			Active:    true,
			Certified: false,
			Official:  false,
			SysRole:   "",
			MemLevel:  0,
			MemExpire: time.Time{},
			Tags:      "",
		},
	}

	// su.SetValidator(map[string]func(string) bool{ })

	lk.Log("%v", user)

	if err := su.ChkInput(user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := su.ChkEmail(user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "waiting verification code in your email")
}

// @Title verify new user's email
// @Summary sign up action, step 2. send back email verification code
// @Description
// @Tags    Sign
// @Accept  multipart/form-data
// @Produce json
// @Param   uname  formData  string  true  "unique user name"
// @Param   code   formData  string  true  "verification code (in user's email)"
// @Success 200 "OK - sign-up successfully"
// @Failure 400 "Fail - incorrect verification code"
// @Failure 500 "Fail - internal error"
// @Router /api/sign/verify-email [post]
func VerifyEmail(c echo.Context) error {
	var (
		uname = c.FormValue("uname")
		code  = c.FormValue("code")
	)

	user, err := su.VerifyCode(uname, code)
	if err != nil || user == nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// double check before storing
	if err := su.ChkInput(user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// store into db
	if err := su.Store(user); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// sign-up ok calling...
	{

	}

	return c.JSON(http.StatusOK, "registered successfully")
}

// @Title sign in
// @Summary sign in action. if ok, got token
// @Description
// @Tags    Sign
// @Accept  multipart/form-data
// @Produce json
// @Param   uname formData string true "user name or email"
// @Param   pwd   formData string true "password" Format(password)
// @Success 200 "OK - sign-in successfully"
// @Failure 400 "Fail - incorrect password"
// @Failure 500 "Fail - internal error"
// @Router /api/sign/in [post]
func LogIn(c echo.Context) error {

	var (
		uname = c.FormValue("uname")
		pwd   = c.FormValue("pwd")
		email = c.FormValue("uname")
	)

	lk.Debug("login: [%v] [%v]", uname, pwd)

	user := &u.User{
		Core: u.Core{
			UName:    uname,
			Password: pwd,
			Email:    email,
		},
		Profile: u.Profile{},
		Admin:   u.Admin{},
	}

AGAIN:

	if err := si.UserStatusIssue(user); err != nil {

		///////////////////////////////////////
		// external user checking
		{
			// try v-site existing user check. if external user already exists, wrap user & login again
			if u := validateSavedExtUser(uname, pwd); u != nil {
				user = u
				goto AGAIN
			}

			// external v-site check via remote api
			if ok, err := vUserLoginCheck(uname, pwd); err == nil && ok {
				// if can login v-site, but doesn't exist, create a new external user, u.uname is like "13888888888@@@V"
				u, err := createExtUser(uname, pwd)
				if err != nil {
					return c.String(http.StatusInternalServerError, "ERR: creating external user, "+err.Error())
				}
				user = u
				goto AGAIN
			}
		}
		///////////////////////////////////////

		return c.String(http.StatusBadRequest, err.Error())
	}

	if !si.PwdOK(user) { // if successful, user updated.
		return c.String(http.StatusBadRequest, "incorrect password")
	}

	// fmt.Println(user)

	// now, user is real user in db
	defer lk.FailOnErr("%v", si.Trail(user.UName)) // Refresh Online Users, here UName is real

	// log in ok calling...
	{
		us, err := fm.UseUser(user.UName)
		if err != nil || us == nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	defer func() { UserCache.Store(user.UName, user) }() // save current user for other usage

	claims := u.MakeClaims(user)
	token := u.GenerateToken(claims)
	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
		"auth":  "Bearer " + token,
	})
}

// @Title reset password
// @Summary reset password action, step 1. send verification code to user's email for authentication
// @Description
// @Tags    Sign
// @Accept  multipart/form-data
// @Produce json
// @Param   uname   formData   string  true  "unique user name"
// @Param   email   formData   string  true  "user's email" Format(email)
// @Success 200 "OK - then waiting for verification code"
// @Failure 400 "Fail - invalid registry fields"
// @Failure 500 "Fail - internal error"
// @Router /api/sign/reset-pwd [post]
func ResetPwd(c echo.Context) error {

	user := &u.User{
		Core: u.Core{
			UName: c.FormValue("uname"),
			Email: c.FormValue("email"),
		},
		Profile: u.Profile{},
		Admin:   u.Admin{},
	}

	if err := rp.UserStatusIssue(user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if !rp.EmailOK(user) {
		return c.String(http.StatusBadRequest, fmt.Sprintf("input email [%s] is different from [%s] sign-up", user.Email, user.UName))
	}

	// load full user before ChkEmail
	user, ok, err := u.LoadUser(user.UName, true)
	if err != nil || !ok {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if err := su.ChkEmail(user); err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "waiting verification code in your email")
}

// @Title update new password
// @Summary reset password action, step 2. send back verification code for updating password
// @Description
// @Tags    Sign
// @Accept  multipart/form-data
// @Produce json
// @Param   uname  formData  string  true  "unique user name"
// @Param   code   formData  string  true  "verification code (in user's email)"
// @Param   pwd    formData  string  true  "new password"
// @Success 200 "OK   - password updated successfully"
// @Failure 400 "Fail - incorrect verification code"
// @Failure 500 "Fail - internal error"
// @Router /api/sign/verify-reset-pwd [post]
func VerifyResetPwd(c echo.Context) error {
	var (
		uname = c.FormValue("uname")
		code  = c.FormValue("code")
		pwd   = c.FormValue("pwd")
	)

	user, err := su.VerifyCode(uname, code)
	if err != nil || user == nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// check new password
	if rst := su.ChkPwd(pwd); rst.OK {
		user.Password = pwd
	} else {
		return c.String(http.StatusBadRequest, "invalid password, at least 11 length with UPPER CASE, number and symbol")
	}

	// store into db
	if err := su.Store(user); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "password updated")
}
