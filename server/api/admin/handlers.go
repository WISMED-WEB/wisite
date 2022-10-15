package admin

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	. "github.com/digisan/go-generics/v2"
	lk "github.com/digisan/logkit"
	u "github.com/digisan/user-mgr/user"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// *** after implementing, register with path in 'admin.go' ***

// @Title get side menu
// @Summary get tailored side menu for different user group
// @Description
// @Tags    Admin
// @Accept  json
// @Produce json
// @Success 200 "OK - get menu successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/admin/spa/menu [get]
// @Security ApiKeyAuth
func Menu(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
	)

	user, ok, err := u.LoadActiveUser(uname)

	switch {
	case err != nil:
		return c.String(http.StatusInternalServerError, err.Error())
	case !ok:
		return c.String(http.StatusInternalServerError, fmt.Sprintf("invalid user status@[%s], dormant?", user.UName))
	}
	// --- //

	var menu []string

	// old version
	// switch user.MemLevel {
	// case 0: // unsubscribe
	// 	menu = []string{"whats-new", "topic", "task"}
	// case 1: // subscribe
	// 	menu = []string{"whats-new", "topic", "bookmark", "my-sharing", "assign", "task", "vote"}
	// case 2: // advanced subscribe
	// 	menu = []string{"whats-new", "topic", "bookmark", "my-sharing", "assign", "task", "vote", "audit"}
	// case 3: // admin
	// 	menu = []string{"whats-new", "topic", "bookmark", "my-sharing", "assign", "task", "vote", "audit", "admin"}
	// default:
	// 	lk.Warn("MemLevel is invalid @ [%v], only 0-3 is valid", user.MemLevel)
	// 	menu = []string{}
	// }

	// simple version
	switch user.MemLevel {
	case 0: // unsubscribe
		menu = []string{"whats-new", "topic"}
	case 1: // subscribe
		menu = []string{"whats-new", "topic", "bookmark", "my-sharing"}
	case 2: // advanced subscribe
		menu = []string{"whats-new", "topic", "bookmark", "my-sharing"}
	case 3: // admin
		menu = []string{"whats-new", "topic", "bookmark", "my-sharing", "admin"}
	default:
		lk.Warn("MemLevel is invalid @ [%v], only 0-3 is valid", user.MemLevel)
		menu = []string{}
	}

	// next version menu ......
	// switch user.MemLevel {
	// case 0:
	// 	menu = []string{"whats-new", "topic", "my-sharing", "my-task"}
	// case 1:
	// 	menu = []string{"whats-new", "topic", "bookmark", "my-sharing", "my-task"}
	// case 2:
	// 	menu = []string{"whats-new", "topic", "bookmark", "my-sharing", "my-task"}
	// case 3:
	// 	menu = []string{"whats-new", "topic", "bookmark", "my-sharing", "my-task", "admin"}
	// default:
	// 	menu = []string{}
	// }

	menu = append(menu, "profile", "wisite-green")

	return c.JSON(http.StatusOK, menu)
}

// @Title list all users
// @Summary get all users' info
// @Description
// @Tags    Admin
// @Accept  json
// @Produce json
// @Param   uname  query string false "user filter with uname wildcard(*)"
// @Param   name   query string false "user filter with name wildcard(*)"
// @Param   active query string false "user filter with active status"
// @Success 200 "OK - list successfully"
// @Failure 401 "Fail - unauthorized error"
// @Failure 500 "Fail - internal error"
// @Router /api/admin/users [get]
// @Security ApiKeyAuth
func ListUser(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
	)

	user, ok, err := u.LoadActiveUser(uname)

	switch {
	case err != nil:
		return c.String(http.StatusInternalServerError, err.Error())
	case !ok:
		return c.String(http.StatusInternalServerError, fmt.Sprintf("invalid user status@[%s], dormant?", user.UName))
	}

	if user.MemLevel != 3 {
		return c.String(http.StatusUnauthorized, "failed, you are not authorized to this api")
	}
	// --- //

	var (
		active = c.QueryParam("active")
		wUname = c.QueryParam("uname")
		wName  = c.QueryParam("name")
		rUname = wc2re(wUname)
		rName  = wc2re(wName)
	)

	users, err := u.ListUser(func(u *u.User) bool {
		switch {
		case len(wUname) > 0 && !rUname.MatchString(u.UName):
			return false
		case len(wName) > 0 && !rName.MatchString(u.Name):
			return false
		case len(active) > 0:
			if bActive, err := strconv.ParseBool(active); err == nil {
				return bActive == u.Active
			}
			return false
		default:
			return true
		}
	})

	for _, user := range users {
		user.Password = strings.Repeat("*", len(user.Password))
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

// @Title get avatar of a user
// @Summary get a user's avatar src as base64
// @Description
// @Tags    Admin
// @Accept  json
// @Produce json
// @Param   uname query string true "uname for its avatar"
// @Success 200 "OK - get avatar src base64"
// @Failure 400 "Fail - user does not exist"
// @Failure 404 "Fail - avatar is empty"
// @Router /api/admin/avatar [get]
// @Security ApiKeyAuth
func UserAvatar(c echo.Context) error {

	uname := c.QueryParam("uname")
	if len(uname) == 0 {
		return c.String(http.StatusBadRequest, uname+" cannot be empty")
	}

	user, ok, err := u.LoadUser(uname, true)
	if err != nil || !ok {
		return c.String(http.StatusBadRequest, "couldn't find user: "+uname)
	}

	atype, b64 := user.AvatarBase64(false)
	if atype == "" || b64 == "" {
		return c.String(http.StatusNotFound, "avatar is empty")
	}

	src := fmt.Sprintf("data:%s;base64,%s", atype, b64)
	return c.JSON(http.StatusOK, struct {
		Src string `json:"src"`
	}{Src: src})
}

// @Title list online users
// @Summary get all online users
// @Description
// @Tags    Admin
// @Accept  json
// @Produce json
// @Param   uname query string false "user filter with uname wildcard(*)"
// @Success 200 "OK - list successfully"
// @Failure 401 "Fail - unauthorized error"
// @Failure 500 "Fail - internal error"
// @Router /api/admin/onlines [get]
// @Security ApiKeyAuth
func ListOnlineUser(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
	)

	user, ok, err := u.LoadActiveUser(uname)

	switch {
	case err != nil:
		return c.String(http.StatusInternalServerError, err.Error())
	case !ok:
		return c.String(http.StatusInternalServerError, fmt.Sprintf("invalid user status@[%s], dormant?", user.UName))
	}

	if user.MemLevel != 3 {
		return c.String(http.StatusUnauthorized, "failed, you are not authorized to this api")
	}
	// --- //

	var (
		wUname = c.QueryParam("uname")
		rUname = wc2re(wUname)
	)

	onlines, err := u.OnlineUsers()
	// for _, user := range onlines {
	// 	fmt.Println(user)
	// }
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	FilterFast(&onlines, func(i int, e *u.UserOnline) bool {
		if len(wUname) > 0 && !rUname.MatchString(e.Uname) {
			return false
		}
		return true
	})

	return c.JSON(http.StatusOK, onlines)
}

// return uname, set flag, return ok, error
func switchField(c echo.Context, fn func(uname string, flag bool) (*u.User, bool, error)) (string, bool, bool, error) {
	uname := c.FormValue("uname")
	flagstr := c.FormValue("flag")
	flag, err := strconv.ParseBool(flagstr)
	if err != nil {
		return "", flag, false, fmt.Errorf("flag must be true/false")
	}
	_, ok, err := fn(uname, flag)
	return uname, flag, ok, err
}

// @Title activate user
// @Summary activate or deactivate a user
// @Description
// @Tags    Admin
// @Accept  multipart/form-data
// @Produce json
// @Param   uname  formData  string  true  "unique user name"
// @Param   flag   formData  string  true  "true: activate, false: deactivate"
// @Success 200 "OK - action successfully"
// @Failure 400 "Fail - invalid true/false flag"
// @Failure 401 "Fail - unauthorized error"
// @Failure 500 "Fail - internal error"
// @Router /api/admin/activate [put]
// @Security ApiKeyAuth
func ActivateUser(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
	)

	user, ok, err := u.LoadActiveUser(uname)

	switch {
	case err != nil:
		return c.String(http.StatusInternalServerError, err.Error())
	case !ok:
		return c.String(http.StatusInternalServerError, fmt.Sprintf("invalid user status@[%s], dormant?", user.UName))
	}

	if user.MemLevel != 3 {
		return c.String(http.StatusUnauthorized, "failed, you are not authorized to this api")
	}
	// --- //

	uname, flag, ok, err := switchField(c, u.ActivateUser)
	if err != nil {
		if uname == "" {
			return c.String(http.StatusBadRequest, err.Error())
		}
		if !ok {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}
	m := map[bool]string{
		true:  "activated",
		false: "deactivated",
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("[%s] is %s", uname, m[flag]))
}

// @Title officialize user
// @Summary officialize or un-officialize a user
// @Description
// @Tags    Admin
// @Accept  multipart/form-data
// @Produce json
// @Param   uname  formData  string  true  "unique user name"
// @Param   flag   formData  string  true  "true: officialize, false: un-officialize"
// @Success 200 "OK - action successfully"
// @Failure 400 "Fail - invalid true/false flag"
// @Failure 401 "Fail - unauthorized error"
// @Failure 500 "Fail - internal error"
// @Router /api/admin/officialize [put]
// @Security ApiKeyAuth
func OfficializeUser(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
	)

	user, ok, err := u.LoadActiveUser(uname)

	switch {
	case err != nil:
		return c.String(http.StatusInternalServerError, err.Error())
	case !ok:
		return c.String(http.StatusInternalServerError, fmt.Sprintf("invalid user status@[%s], dormant?", user.UName))
	}

	if user.MemLevel != 3 {
		return c.String(http.StatusUnauthorized, "failed, you are not authorized to this api")
	}
	// --- //

	uname, flag, ok, err := switchField(c, u.OfficializeUser)
	if err != nil {
		if uname == "" {
			return c.String(http.StatusBadRequest, err.Error())
		}
		if !ok {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}
	m := map[bool]string{
		true:  "switched to official account",
		false: "switched to unofficial account",
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("[%s] is %s", uname, m[flag]))
}
