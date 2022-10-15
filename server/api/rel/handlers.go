package rel

import (
	"fmt"
	"net/http"

	lk "github.com/digisan/logkit"
	r "github.com/digisan/user-mgr/relation"
	u "github.com/digisan/user-mgr/user"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// *** after implementing, register with path in 'rel.go' *** //

// @Title   user relations
// @Summary relation actions
// @Description
// @Tags    Relation
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
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
	)

	if _, ok, err := u.LoadUser(uname, true); err != nil || !ok {
		return c.String(http.StatusInternalServerError, "couldn't find user: "+uname)
	}

	var (
		whom     = c.Param("whom")
		action   = c.QueryParam("action")
		mActFlag = map[string]int{
			"follow":   r.FOLLOW,
			"unfollow": r.UNFOLLOW,
			"block":    r.BLOCK,
			"unblock":  r.UNBLOCK,
			"mute":     r.MUTE,
			"unmute":   r.UNMUTE,
		}
	)

	flag, ok := mActFlag[action]
	if !ok {
		return c.String(http.StatusBadRequest, "invalid action, only accept [follow, unfollow, block, unblock, mute, unmute]")
	}

	if err := r.RelAction(uname, flag, whom); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("%s %s successfully now", action, whom))
}

// @Title user relation content
// @Summary get all relation users for one type
// @Description
// @Tags    Relation
// @Accept  json
// @Produce json
// @Param   type path string true "relation content type to apply, accept [following, follower, blocked, muted]"
// @Success 200 "OK - got following successfully"
// @Failure 400 "Fail - invalid relation content type"
// @Failure 500 "Fail - internal error"
// @Router /api/rel/content/{type} [get]
// @Security ApiKeyAuth
func GetContent(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
	)

	if _, ok, err := u.LoadUser(uname, true); err != nil || !ok {
		return c.String(http.StatusInternalServerError, "couldn't find user: "+uname)
	}

	var (
		contType  = c.Param("type")
		mContFlag = map[string]int{
			"following": r.FOLLOWING,
			"follower":  r.FOLLOWER,
			"blocked":   r.BLOCKED,
			"muted":     r.MUTED,
		}
	)

	lk.Debug("%s", contType)

	flag, ok := mContFlag[contType]
	if !ok {
		return c.String(http.StatusBadRequest, "invalid action, only accept [following, follower, blocked, muted]")
	}

	names, err := r.ListRel(uname, flag, true)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, names)
}
