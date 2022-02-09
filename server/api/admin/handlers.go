package admin

import (
	"fmt"
	"net/http"

	"github.com/digisan/user-mgr/udb"
	usr "github.com/digisan/user-mgr/user"
	"github.com/labstack/echo/v4"
)

// *** after implementing, register with path in 'api_reg.go' ***

// SignIn godoc
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
		return c.String(http.StatusInternalServerError, fmt.Sprintln(err))
	}
	return c.JSON(http.StatusOK, users)
}

// SignIn godoc
// @Title list online users
// @Summary get all online users' info in db
// @Description
// @Tags    admin
// @Accept  json
// @Produce json
// @Success 200 "OK - list successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/admin/onlineusers [get]
func ListOnlineUser(c echo.Context) error {
	users, err := udb.UserDB.ListOnlineUsers()
	// for _, user := range users {
	// 	fmt.Println(user)
	// }
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintln(err))
	}
	return c.JSON(http.StatusOK, users)
}
