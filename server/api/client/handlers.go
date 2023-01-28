package client

import (
	"fmt"
	"net/http"

	lk "github.com/digisan/logkit"
	u "github.com/digisan/user-mgr/user"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// @Title set client browser's viewport
// @Summary set client browser's viewport ( width, height )
// @Description
// @Tags    Client
// @Accept  json
// @Produce json
// @Param   innerSize body string true "width: window.innerWidth & height: window.innerHeight"
// @Success 200 "OK - set client viewport ok"
// @Failure 400 "Fail - invalid width or height for setting viewport"
// @Router /api/client/set/view [put]
// @Security ApiKeyAuth
func SetClientView(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
		cv      = new(Area)
	)
	if err := c.Bind(cv); err != nil || cv.Width <= 0 || cv.Height <= 0 {
		return c.String(http.StatusBadRequest, "set client viewport error")
	}
	AddLayout(uname, newLayout(cv))

	mLayout.Range(func(key, value any) bool {
		lk.Log("%v: %v", key, value)
		return true
	})

	return c.JSON(http.StatusOK, "Set ClientView Successfully")
}

// @Title get client viewport & others' size
// @Summary get client viewport, header, menu, content, & post-title size
// @Description
// @Tags    Client
// @Accept  json
// @Produce json
// @Success 200 "OK - get client viewport & other parts' size ok"
// @Failure 400 "Fail - viewport is not set"
// @Router /api/client/get/size [get]
// @Security ApiKeyAuth
func GetSize(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
		lo      = GetLayout(uname)
	)
	if lo == nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("[%v]'s viewport is NOT set", uname))
	}

	type PostArea struct {
		Title   Area `json:"title"`
		Content Area `json:"content"`
	}

	return c.JSON(http.StatusOK, struct {
		Layout   Layout   `json:"layout"`
		PostArea PostArea `json:"postarea"`
	}{
		Layout: *lo,
		PostArea: PostArea{
			Title: Area{
				Width:  lo.PostWidth(),
				Height: lo.PostTitleHeight(),
			},
			Content: Area{
				Width:  lo.PostWidth(),
				Height: lo.PostContentHeight(),
			},
		},
	})
}
