package user

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	su "github.com/digisan/user-mgr/sign-up"
	"github.com/digisan/user-mgr/udb"
	usr "github.com/digisan/user-mgr/user"
	vf "github.com/digisan/user-mgr/user/valfield"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// @Title get user profile
// @Summary get user profile
// @Description
// @Tags    user
// @Accept  json
// @Produce json
// @Success 200 "OK - profile get successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/user/profile [get]
// @Security ApiKeyAuth
func Profile(c echo.Context) error {
	userTkn := c.Get("user").(*jwt.Token)
	claims := userTkn.Claims.(*usr.UserClaims)
	uname := claims.UName

	u, ok, err := udb.UserDB.LoadUser(uname, true)
	if err != nil || !ok {
		return c.String(http.StatusInternalServerError, "Couldn't find user: "+uname)
	}
	return c.JSON(http.StatusOK, *u)
}

// @Title set user profile
// @Summary set user profile
// @Description
// @Tags    user
// @Accept  multipart/form-data
// @Produce json
// @Param   phone     formData   string  false  "phone number"
// @Param   addr      formData   string  false  "address"
// @Param   nid       formData   string  false  "national ID"
// @Param   gender    formData   string  false  "gender M/F"
// @Param   position  formData   string  false  "job position"
// @Param   title     formData   string  false  "title"
// @Param   employer  formData   string  false  "employer"
// @Param   avatar    formData   file    false  "avatar"
// @Success 200 "OK - profile set successfully"
// @Failure 400 "Fail - invalid set fields"
// @Failure 500 "Fail - internal error"
// @Router /api/user/setprofile [post]
// @Security ApiKeyAuth
func SetProfile(c echo.Context) error {
	userTkn := c.Get("user").(*jwt.Token)
	claims := userTkn.Claims.(*usr.UserClaims)
	uname := claims.UName

	u, ok, err := udb.UserDB.LoadUser(uname, true)
	if err != nil || !ok {
		return c.String(http.StatusInternalServerError, "Couldn't find user: "+uname)
	}

	u.Phone = c.FormValue("phone")
	u.Addr = c.FormValue("addr")
	u.NationalID = c.FormValue("nid")
	u.Gender = c.FormValue("gender")
	u.Position = c.FormValue("position")
	u.Title = c.FormValue("title")
	u.Employer = c.FormValue("employer")

	// Read & Set Avatar
	file, err := c.FormFile("avatar")
	var ext string
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			goto VALIDATE
		}
		return c.String(http.StatusBadRequest, err.Error())
	}
	ext = strings.TrimPrefix(filepath.Ext(file.Filename), ".")
	if err := u.SetAvatarByFormFile("image/"+ext, file); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

VALIDATE:
	// validate
	if err := su.ChkInput(u, vf.UName, vf.SysRole, vf.MemLevel, vf.MemExpire, vf.Tags); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// update
	if err := udb.UserDB.UpdateUser(u); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "Profile Updated")
}

// @Title get user avatar
// @Summary get user avatar src as base64
// @Description
// @Tags    user
// @Accept  json
// @Produce json
// @Success 200 "OK - get avatar src base64"
// @Failure 404 "Fail - avatar is empty"
// @Failure 500 "Fail - internal error"
// @Router /api/user/avatar [get]
// @Security ApiKeyAuth
func Avatar(c echo.Context) error {
	userTkn := c.Get("user").(*jwt.Token)
	claims := userTkn.Claims.(*usr.UserClaims)
	uname := claims.UName

	u, ok, err := udb.UserDB.LoadUser(uname, true)
	if err != nil || !ok {
		return c.String(http.StatusInternalServerError, "Couldn't find user: "+uname)
	}

	atype, b64 := u.AvatarBase64(false)
	if atype == "" || b64 == "" {
		return c.String(http.StatusNotFound, "avatar is empty")
	}

	src := fmt.Sprintf("data:%s;base64,%s", atype, b64)
	return c.JSON(http.StatusOK, struct {
		Src string `json:"src"`
	}{Src: src})
}

// u := &usr.User{
// 	Active:     "T",
// 	UName:      c.FormValue("uname"),
// 	Email:      c.FormValue("email"),
// 	Name:       c.FormValue("name"),
// 	Password:   c.FormValue("pwd"),
// 	Regtime:    "TBD",
// 	Phone:      "",      //
// 	Addr:       "",      //
// 	SysRole:    "",   **
// 	MemLevel:   "0",  **
// 	MemExpire:  "",   **
// 	NationalID: "",      //
// 	Gender:     "",      //
// 	Position:   "",      //
// 	Title:      "",      //
// 	Employer:   "",      //
// 	Tags:       "",   **
// 	Avatar:     []byte{}, //
// }
