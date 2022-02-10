package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/digisan/user-mgr/udb"
	usr "github.com/digisan/user-mgr/user"
	"github.com/labstack/echo/v4"
)

// *** after implementing, register with path in 'admin.go' ***

// @Title list all users
// @Summary get all users' info in db
// @Description
// @Tags    admin
// @Accept  json
// @Produce json
// @Success 200 "OK - list successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/admin/users [get]
func ListUser(c echo.Context) error {
	users, err := udb.UserDB.ListUsers(func(u *usr.User) bool {
		return true
	})
	// for _, user := range users {
	// 	fmt.Println(user)
	// }
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

// @Title list online users
// @Summary get all online users' info in db
// @Description
// @Tags    admin
// @Accept  json
// @Produce json
// @Success 200 "OK - list successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/admin/onlines [get]
func ListOnlineUser(c echo.Context) error {
	users, err := udb.UserDB.ListOnlineUsers()
	// for _, user := range users {
	// 	fmt.Println(user)
	// }
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

// @Title activate user
// @Summary activate or deactivate a user
// @Description
// @Tags    admin
// @Accept  multipart/form-data
// @Produce json
// @Param   uname  formData  string  true  "unique user name"
// @Param   flag   formData  string  true  "true: activate, false: deactivate"
// @Success 200 "OK - action successfully"
// @Failure 400 "Fail - invalid uname"
// @Failure 500 "Fail - internal error"
// @Router /api/admin/activate [post]
func ActivateUser(c echo.Context) error {
	uname := c.FormValue("uname")
	flagstr := c.FormValue("flag")
	flag, err := strconv.ParseBool(flagstr)
	if err != nil {
		return c.String(http.StatusBadRequest, "flag must be true/false")
	}
	_, err = udb.UserDB.ActivateUser(uname, flag)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	m := map[bool]string{
		true:  "activated",
		false: "deactivated",
	}
	return c.String(http.StatusOK, fmt.Sprintf("[%s] is %s", uname, m[flag]))
}
