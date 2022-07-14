package post

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	em "github.com/digisan/event-mgr"
	. "github.com/digisan/go-generics/v2"
	fd "github.com/digisan/gotk/filedir"
	gio "github.com/digisan/gotk/io"
	lk "github.com/digisan/logkit"
	u "github.com/digisan/user-mgr/user"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// https://github.com/swaggo/swag
// https://swagger.io/specification/

// *** after implementing, register with path in 'post.go' *** //

// @Title Post template
// @Summary get Post template for dev reference.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Success 200 "OK - upload successfully"
// @Router /api/post/template [get]
// @Security ApiKeyAuth
func Template(c echo.Context) error {
	return c.JSON(http.StatusOK, Post{
		Category: "post category",
		Topic:    "post topic",
		Content: []Paragraph{
			{
				Text: "some words for this attach",
				Path: "attached stuff path, which should have been given from 'file upload'",
			},
		},
		Summary: "summarize your topic",
	})
}

// @Title upload a Post
// @Summary upload a Post by filling a Post template.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Param   data body string true "filled Post template json file"
// @Param   followee  query string false "followee Post ID (empty when doing a new post)"
// @Success 200 "OK - upload successfully"
// @Failure 400 "Fail - incorrect Post format"
// @Failure 500 "Fail - internal error"
// @Router /api/post/upload [post]
// @Security ApiKeyAuth
func Upload(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
		flwee   = c.QueryParam("followee")
	)

	P := new(Post)
	if err := c.Bind(P); err != nil {
		return c.String(http.StatusBadRequest, "incorrect Post format: "+err.Error())
	}
	// lk.Log("---> %s -- %v", uname, P)

	// get rid of empty paragraph
	//
	FilterFast(&P.Content, func(i int, e Paragraph) bool {
		return len(e.Text) > 0 || len(e.Path) > 0
	})

	// validate each path from P
	//
	paths := []string{}
	for _, item := range P.Content {
		if len(item.Path) > 0 {
			path := filepath.Join("data/user-space", uname, item.Path)
			paths = append(paths, path)
		}
	}
	ok, epath := fd.AllExistAsWhole(paths...)
	if !ok {
		return c.String(http.StatusBadRequest, fmt.Sprintf("'%s' is invalid storage at server", filepath.Base(epath)))
	}

	// save P as JSON for event
	//
	data, err := json.Marshal(P)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	evt := em.NewEvent("", uname, "Post", string(data))
	if len(evt.ID) > 0 {
		if err = em.AddEvent(evt); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// DEBUG
		gio.MustAppendFile("./debug.txt", []byte(evt.ID), true)

		// FOLLOWING...
		if len(flwee) > 0 {
			ef, err := em.FetchFollow(flwee)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
			if ef == nil {
				ef = em.NewEventFollow(flwee)
			}
			if err := ef.AddFollower(evt.ID); err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
		}
	}

	// lk.Log("---> %s", em.CurrIDs())

	return c.JSON(http.StatusOK, evt)
}

// @Title get a batch of Post id group
// @Summary get a batch of Post id group.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Param   fetchby query string true "time or count"
// @Param   value   query string true "recent [value] minutes for time OR most recent [value] count"
// @Success 200 "OK - get successfully"
// @Failure 400 "Fail - incorrect query param type"
// @Failure 404 "Fail - not found"
// @Failure 500 "Fail - internal error"
// @Router /api/post/ids [get]
// @Security ApiKeyAuth
func IdBatch(c echo.Context) error {
	var (
		fetchby = c.QueryParam("fetchby")
		value   = c.QueryParam("value")
		ids     = []string{}
	)

	if fetchby = strings.ToLower(fetchby); NotIn(fetchby, "time", "count") {
		return c.String(http.StatusBadRequest, "'fetchby' must be one of [time, count]")
	}

	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "'value' must be a valid number for time(minutes) or count")
	}

	switch fetchby {
	case "time":
		ids, err = em.FetchEvtIDsByTm(value + "m")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	case "count":
		ids, err = em.FetchEvtIDsByCnt(int(n), "")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	// lk.Log("IdBatch ---> %d : %v", len(ids), ids)

	// if len(ids) == 0 {
	// 	return c.JSON(http.StatusNotFound, ids)
	// }
	return c.JSON(http.StatusOK, ids)
}

// @Title get all Post id group
// @Summary get all Post id group.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Success 200 "OK - get successfully"
// @Failure 404 "Fail - empty event ids"
// @Failure 500 "Fail - internal error"
// @Router /api/post/ids-all [get]
// @Security ApiKeyAuth
func IdAll(c echo.Context) error {

	ids, err := em.FetchEvtIDs(nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// lk.Log("IdAll ---> %d : %v", len(ids), ids)

	// if len(ids) == 0 {
	// 	return c.JSON(http.StatusNotFound, ids)
	// }
	return c.JSON(http.StatusOK, ids)
}

// @Title get one Post content
// @Summary get one Post content.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Param   id   query string true "Post ID for its content"
// @Success 200 "OK - get successfully"
// @Failure 400 "Fail - incorrect query param id"
// @Failure 404 "Fail - not found"
// @Failure 500 "Fail - internal error"
// @Router /api/post/one [get]
// @Security ApiKeyAuth
func GetOne(c echo.Context) error {
	var (
		id = c.QueryParam("id")
	)

	if len(id) == 0 {
		c.String(http.StatusBadRequest, "'id' is invalid (cannot be empty)")
	}

	content, err := em.FetchEvent(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	if content == nil {
		c.String(http.StatusNotFound, fmt.Sprintf("Post not found @%s", id))
	}
	return c.JSON(http.StatusOK, content)
}

// @Title get own Post id group in a specific period
// @Summary get own Post id group in one specific time period.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Param   period query string false "time period for query, format is 'yyyymm', e.g. '202206'. if missing, current yyyymm applies"
// @Success 200 "OK - get successfully"
// @Failure 400 "Fail - incorrect query param type"
// @Failure 404 "Fail - empty event ids"
// @Failure 500 "Fail - internal error"
// @Router /api/post/own/ids [get]
// @Security ApiKeyAuth
func IdOwn(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
		period  = c.QueryParam("period")
	)

	if len(period) == 0 {
		period = time.Now().Format("200601")
	}
	if _, err := time.Parse("200601", period); err != nil {
		return c.String(http.StatusBadRequest, "'period' format must be 'yyyymm', e.g. '202206'")
	}

	lk.Log("%s -- %s", uname, period)

	ids, err := em.FetchOwn(uname, period)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	// if len(ids) == 0 {
	// 	return c.JSON(http.StatusNotFound, ids)
	// }
	return c.JSON(http.StatusOK, ids)
}

// @Title get a Post follower-Post ids
// @Summary get a specified Post follower-Post id group.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Param   followee query string true "followee Post ID"
// @Success 200 "OK - get successfully"
// @Failure 404 "Fail - empty follower ids"
// @Failure 500 "Fail - internal error"
// @Router /api/post/follower/ids [get]
// @Security ApiKeyAuth
func Followers(c echo.Context) error {
	var (
		flwee = c.QueryParam("followee")
	)

	flwers, err := em.Followers(flwee)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	// if len(flwers) == 0 {
	// 	return c.JSON(http.StatusNotFound, flwers)
	// }
	return c.JSON(http.StatusOK, flwers)
}
