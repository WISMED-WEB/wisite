package user

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	lk "github.com/digisan/logkit"
	si "github.com/digisan/user-mgr/sign-in"
	su "github.com/digisan/user-mgr/sign-up"
	u "github.com/digisan/user-mgr/user"
	vf "github.com/digisan/user-mgr/user/valfield"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// @Title user heartbeats
// @Summary frequently call this to indicate that front-end user is active.
// @Description
// @Tags    User
// @Accept  json
// @Produce json
// @Success 200 "OK - heartbeats successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/user/heartbeats [patch]
// @Security ApiKeyAuth
func HeartBeats(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
	)

	if err := si.Trail(uname); err != nil {
		lk.Debug("%v", err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, fmt.Sprintf("[%v] heartbeats", uname))
}

// @Title get user profile
// @Summary get user profile
// @Description
// @Tags    User
// @Accept  json
// @Produce json
// @Success 200 "OK - profile get successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/user/profile [get]
// @Security ApiKeyAuth
func Profile(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
	)

	user, ok, err := u.LoadUser(uname, true)
	if err != nil || !ok {
		return c.String(http.StatusInternalServerError, "couldn't find user: "+uname)
	}

	if len(user.Profile.Avatar) > 32 {
		user.Profile.Avatar = user.Profile.Avatar[:32]
	}

	return c.JSON(http.StatusOK, struct {
		u.Profile
		Uname      string `json:"uname"`
		Email      string `json:"email"`
		MemberDays string `json:"memberDays"`
	}{
		user.Profile,
		user.UName,
		user.Email,
		fmt.Sprintf("%v", int(user.SinceJoined().Hours()/24.0)),
	})
}

// @Title set user profile
// @Summary set user profile
// @Description
// @Tags    User
// @Accept  multipart/form-data
// @Produce json
// @Param   phone     formData   string  false  "phone number"
// @Param   addr      formData   string  false  "address"
// @Param   city      formData   string  false  "city"
// @Param   country   formData   string  false  "country"
// @Param   pidtype   formData   string  false  "personal id type"
// @Param   pid       formData   string  false  "personal id"
// @Param   gender    formData   string  false  "gender M/F"
// @Param   dob       formData   string  false  "date of birth"
// @Param   position  formData   string  false  "job position"
// @Param   title     formData   string  false  "title"
// @Param   employer  formData   string  false  "employer"
// @Param   bio       formData   string  false  "biography"
// @Param   avatar    formData   file    false  "avatar"
// @Success 200 "OK - profile set successfully"
// @Failure 400 "Fail - invalid set fields"
// @Failure 500 "Fail - internal error"
// @Router /api/user/setprofile [post]
// @Security ApiKeyAuth
func SetProfile(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
	)

	user, ok, err := u.LoadActiveUser(uname)
	if err != nil || !ok {
		return c.String(http.StatusInternalServerError, "couldn't find user: "+uname)
	}

	user.Phone = c.FormValue("phone")
	user.Addr = c.FormValue("addr")
	user.City = c.FormValue("city")
	user.Country = c.FormValue("country")
	user.PersonalIDType = c.FormValue("pidtype")
	user.PersonalID = c.FormValue("pid")
	user.Gender = c.FormValue("gender")
	user.DOB = c.FormValue("dob")
	user.Position = c.FormValue("position")
	user.Title = c.FormValue("title")
	user.Employer = c.FormValue("employer")
	user.Bio = c.FormValue("bio")

	// Read & Set Avatar
	file, err := c.FormFile("avatar")
	var ext string
	if err != nil && file == nil {
		e := err.Error()
		if strings.Contains(e, "no such file") || strings.Contains(e, "no multipart boundary param in Content-Type") {
			goto VALIDATE // if no file submitted, do nothing
		}
		return c.String(http.StatusBadRequest, e)
	}
	ext = strings.TrimPrefix(filepath.Ext(file.Filename), ".")
	if err := user.SetAvatarByFormFile("image/"+ext, file); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

VALIDATE:
	// validate
	if err := su.ChkInput(user, vf.UName, vf.EmailDB, vf.SysRole, vf.MemLevel, vf.MemExpire, vf.Tags); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// update
	if err := u.UpdateUser(user); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "Profile Updated")
}

// @Title get self avatar
// @Summary get self avatar src as base64
// @Description
// @Tags    User
// @Accept  json
// @Produce json
// @Success 200 "OK - get avatar src base64"
// @Failure 404 "Fail - avatar is empty"
// @Failure 500 "Fail - internal error"
// @Router /api/user/avatar [get]
// @Security ApiKeyAuth
func Avatar(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
	)

	user, ok, err := u.LoadUser(uname, true)
	if err != nil || !ok {
		return c.String(http.StatusInternalServerError, "couldn't find user: "+uname)
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

// var u = &u.User{
// 	Core: u.Core{
// 		UName:    "",
// 		Email:    "",
// 		Password: "",
// 	},
// 	Profile: u.Profile{
// 		Name:           "",
// 		Phone:          "",
// 		Country:        "",
// 		City:           "",
// 		Addr:           "",
// 		PersonalIDType: "",
// 		PersonalID:     "",
// 		Gender:         "",
// 		DOB:            "",
// 		Position:       "",
// 		Title:          "",
// 		Employer:       "",
// 		Bio:            "",
// 		AvatarType:     "",
// 		Avatar:         []byte{},
// 	},
// 	Admin: u.Admin{
// 		Regtime:   time.Now().Truncate(time.Second),
// 		Active:    true,
// 		Certified: false,
// 		Official:  false,
// 		SysRole:   "",
// 		MemLevel:  "",
// 		MemExpire: time.Time{},
// 		Tags:      "",
// 	},
// }
