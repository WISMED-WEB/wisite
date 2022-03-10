package rel

import (
	"fmt"
	"net/http"

	rel "github.com/digisan/user-mgr/relation"
	. "github.com/digisan/user-mgr/relation/enum"
	"github.com/digisan/user-mgr/udb"
	usr "github.com/digisan/user-mgr/user"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// *** after implementing, register with path in 'rel.go' *** //

// @Title   user relations
// @Summary relation actions
// @Description
// @Tags    relation
// @Accept  json
// @Produce json
// @Param   action query string true "which action to apply, accept [follow, unfollow, block, unblock, mute, unmute]"
// @Param   whom path string true "whose uname you want to follow"
// @Success 200 "OK - following successfully"
// @Failure 400 "Fail - invalid action type"
// @Failure 500 "Fail - internal error"
// @Router /api/rel/action/{whom} [put]
// @Security ApiKeyAuth
func Action(c echo.Context) error {
	userTkn := c.Get("user").(*jwt.Token)
	claims := userTkn.Claims.(*usr.UserClaims)
	uname := claims.UName

	_, ok, err := udb.UserDB.LoadUser(uname, true)
	if err != nil || !ok {
		return c.String(http.StatusInternalServerError, "couldn't find user: "+uname)
	}

	var (
		whom     = c.Param("whom")
		action   = c.QueryParam("action")
		mActFlag = map[string]int{
			"follow":   DO_FOLLOW,
			"unfollow": DO_UNFOLLOW,
			"block":    DO_BLOCK,
			"unblock":  DO_UNBLOCK,
			"mute":     DO_MUTE,
			"unmute":   DO_UNMUTE,
		}
	)

	flag, ok := mActFlag[action]
	if !ok {
		return c.String(http.StatusBadRequest, "invalid action value, accept [follow, unfollow, block, unblock, mute, unmute]")
	}

	if err := rel.RelAction(uname, flag, whom); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, fmt.Sprintf("%s %s successfully now", action, whom))
}

// @Title user relation content
// @Summary get all relation users for one type
// @Description
// @Tags    relation
// @Accept  json
// @Produce json
// @Param   type path string true "relation content type to apply, accept [following, follower, blocked, muted]"
// @Success 200 "OK - got following successfully"
// @Failure 400 "Fail - invalid relation content type"
// @Failure 500 "Fail - internal error"
// @Router /api/rel/content/{type} [get]
// @Security ApiKeyAuth
func GetContent(c echo.Context) error {
	userTkn := c.Get("user").(*jwt.Token)
	claims := userTkn.Claims.(*usr.UserClaims)
	uname := claims.UName

	_, ok, err := udb.UserDB.LoadUser(uname, true)
	if err != nil || !ok {
		return c.String(http.StatusInternalServerError, "couldn't find user: "+uname)
	}

	var (
		contType  = c.QueryParam("type")
		mContFlag = map[string]int{
			"following": FOLLOWING,
			"follower":  FOLLOWER,
			"blocked":   BLOCKED,
			"muted":     MUTED,
		}
	)

	flag, ok := mContFlag[contType]
	if !ok {
		return c.String(http.StatusBadRequest, "invalid action value, accept [following, follower, blocked, muted]")
	}

	return c.JSON(http.StatusOK, rel.RelContent(uname, flag))
}
